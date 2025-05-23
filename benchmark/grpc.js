import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';
import { CONFIG } from './config.js';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export let options = {
    stages: CONFIG.stages,
    thresholds: {
        // gRPC 性能指标
        'grpc_req_duration': [
            'p(95)<500',    // 95%的请求应该在500ms内完成
            'p(90)<400',    // 90%的请求应该在400ms内完成
            'max<2000'      // 最大响应时间不超过2s
        ],
        'checks': [
            'rate>0.99'     // 检查通过率应该大于99%
        ],
        'vus': [
            'value>0'       // 确保虚拟用户数被记录
        ],
        'vus_max': [
            'value>0'       // 确保最大虚拟用户数被记录
        ],
        'iterations': [
            'rate>0'        // 确保迭代率被记录
        ],
        'data_sent': [
            'rate>0'        // 确保发送数据率被记录
        ],
        'data_received': [
            'rate>0'        // 确保接收数据率被记录
        ],
        'iteration_duration': [
            'p(95)<1500'    // 95%的迭代应该在1.5s内完成
        ]
    },
};

const client = new grpc.Client();
// console.log('Loading proto file from:', CONFIG.grpc.protoFile);

try {
    // 先importPaths，再protoFiles, 引入要测试的proto文件
    client.load([CONFIG.grpc.protoDir], CONFIG.grpc.protoFile);
} catch (e) {
    console.error('Error loading proto file:', e);
}

// TODO 如果需要鉴权，修改这里, 根据不同的场景做不同的测试
export default function () {
    // console.log('Connecting to gRPC server at:', CONFIG.grpc.baseUrl);
    
    // 连接gRPC服务，添加重试和超时配置
    try {
        client.connect(CONFIG.grpc.baseUrl, {
            plaintext: true,
            timeout: '5s'  // 连接超时时间
        });

        // 记录连接成功
        // console.log('gRPC connection successful');

        // 调用gRPC服务
        const userResponse = client.invoke(CONFIG.grpc.method, CONFIG.grpc.methodParams, {
            metadata: CONFIG.grpc.methodMetadata,
            timeout: '10s'  // 调用超时时间
        });
        
        // 所有的检测都通过才算成功, 1. 状态， 2. 响应时间， 3. 消息
        check(userResponse, {
            'get user info status is OK': (r) => r && r.status === grpc.StatusOK,
            'error is null': (r) => r && r.error === null,
            'message info': (r) => {
                // 如果需要这里可以用来判断真实请求后返回的数据的正确性
                // console.log('Response:', r);
                return true;
            }
        });

    } catch (e) {
        // 记录连接失败
        check(null, {
            'gRPC connection successful': () => false,
            'connection error details': () => {
                console.error('gRPC connection error:', e);
                return false;
            }
        });
    } finally {
        // 确保连接被关闭
        try {
            client.close();
        } catch (e) {
            console.error('Error closing gRPC connection:', e);
        }
    }

    sleep(1);
}

export function handleSummary(data) {
    // console.log('data', data);
    
    // 获取请求统计
    const totalRequests = data.metrics.iterations.values.count;
    const failedChecks = data.metrics.checks.values.fails;
    const successRate = ((totalRequests - failedChecks) / totalRequests * 100).toFixed(2);
    
    // 获取最大并发量
    const maxVUs = data.metrics.vus_max.values.max;
    
    // 获取响应时间统计
    const avgDuration = data.metrics.grpc_req_duration.values.avg;
    const maxDuration = data.metrics.grpc_req_duration.values.max;
    const p90Duration = data.metrics.grpc_req_duration.values['p(90)'];
    const p95Duration = data.metrics.grpc_req_duration.values['p(95)'];
    
    // 获取数据统计
    const dataSent = data.metrics.data_sent.values;
    const dataReceived = data.metrics.data_received.values;
    
    // 在控制台输出统计信息
    console.log(`\n=== gRPC 测试统计 ===`);
    console.log(`总请求数: ${totalRequests}`);
    console.log(`检查失败数: ${failedChecks}`);
    console.log(`成功率: ${successRate}%`);
    console.log(`最大并发量: ${maxVUs}`);
    console.log(`\n=== 响应时间统计 ===`);
    console.log(`平均响应时间: ${avgDuration.toFixed(2)}ms`);
    console.log(`最大响应时间: ${maxDuration.toFixed(2)}ms`);
    console.log(`90%响应时间: ${p90Duration.toFixed(2)}ms`);
    console.log(`95%响应时间: ${p95Duration.toFixed(2)}ms`);
    console.log(`\n=== 数据传输统计 ===`);
    console.log(`发送数据: ${(dataSent.count/1024).toFixed(2)}KB (${(dataSent.rate/1024).toFixed(2)}KB/s)`);
    console.log(`接收数据: ${(dataReceived.count/1024).toFixed(2)}KB (${(dataReceived.rate/1024).toFixed(2)}KB/s)`);
    console.log(`====================\n`);

    // 构建统计数据的 JSON 对象
    const statsJson = {
        "测试基本信息": {
            "总请求数": totalRequests,
            "检查失败数": failedChecks,
            "成功率": successRate + "%",
            "最大并发用户数": maxVUs
        },
        "响应时间统计": {
            "平均响应时间": avgDuration.toFixed(2) + "ms",
            "最大响应时间": maxDuration.toFixed(2) + "ms",
            "90%响应时间": p90Duration.toFixed(2) + "ms",
            "95%响应时间": p95Duration.toFixed(2) + "ms"
        },
        "数据传输统计": {
            "发送数据": {
                "总量": (dataSent.count/1024).toFixed(2) + "KB",
                "速率": (dataSent.rate/1024).toFixed(2) + "KB/s"
            },
            "接收数据": {
                "总量": (dataReceived.count/1024).toFixed(2) + "KB",
                "速率": (dataReceived.rate/1024).toFixed(2) + "KB/s"
            }
        }
    };

    return {
        "reports/grpc-report.html": htmlReport(data, {
            title: "gRPC 性能测试报告",
            json: true,
            includeMetrics: true,
            includeThresholds: true,
            includeGroups: true,
            includeChecks: true,
            includeTags: true
        }),
        "reports/grpc-stats.json": JSON.stringify(statsJson, null, 2)
    };
} 