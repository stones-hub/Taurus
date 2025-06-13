package controller

import (
	"encoding/base64"
	"net/http"
	"os"
)

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

func ServeMarkdownDoc(w http.ResponseWriter, r *http.Request) {
	// 读取Markdown文件
	mdContent, err := os.ReadFile("./README.md")
	if err != nil {
		http.Error(w, "无法读取文档", http.StatusInternalServerError)
		return
	}

	// 编码Markdown内容为Base64
	encodedContent := base64.StdEncoding.EncodeToString(mdContent)

	html := `<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>API 接口文档</title>
  <link rel="stylesheet" href="/static/css/github-markdown.min.css">
  <style>
    .markdown-body {
      box-sizing: border-box;
      min-width: 200px;
      max-width: 980px;
      margin: 0 auto;
      padding: 45px;
    }
    .toc {
      position: fixed;
      left: 20px;
      top: 20px;
      padding: 20px;
      width: 280px;
      max-height: 80vh;
      overflow-y: auto;
      background-color: #f8f8f8;
      border: 1px solid #e8e8e8;
      border-radius: 5px;
    }
    .toc-list {
      list-style: none;
      padding-left: 15px;
    }
    .toc-list-item {
      margin: 5px 0;
    }
    .toc-list-item a {
      color: #0366d6;
      text-decoration: none;
      font-family: 'Arial', sans-serif;
      transition: color 0.3s ease;
    }
    .toc-list-item a:hover {
      text-decoration: underline;
      color: #0056b3;
    }
    .toc-list-item.open > a {
      color: #ff0000; /* 更显眼的红色 */
    }
    .toc-list-item {
      cursor: pointer;
    }
    .toc-list-item ul {
      display: none;
    }
    .toc-list-item.open > ul {
      display: block;
    }
    .content-with-toc {
      margin-left: 320px;
    }
    @media (max-width: 1300px) {
      .toc {
        display: none;
      }
      .content-with-toc {
        margin-left: auto;
      }
    }
  </style>
</head>
<body>
  <div id="toc" class="toc"></div>
  <div id="content" class="markdown-body content-with-toc"></div>
  <script src="/static/js/marked.min.js"></script>
  <script>
    // 解码Base64字符串为UTF-8文本
    function decodeBase64UTF8(base64) {
      try {
        const binaryString = atob(base64);
        const bytes = new Uint8Array(binaryString.length);
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i);
        }
        return new TextDecoder('utf-8').decode(bytes);
      } catch (e) {
        console.error('解码失败:', e);
        return 'Error decoding content';
      }
    }
    
    // 创建目录
    function generateTOC(markdown) {
      const headings = [];
      const lines = markdown.split('\n');
      
      // 找出所有标题
      for(let i = 0; i < lines.length; i++) {
        const line = lines[i].trim();
        if(line.startsWith('# ')) {
          headings.push({ level: 1, text: line.substring(2), anchor: 'heading-' + i });
          lines[i] = '<h1 id="heading-' + i + '">' + line.substring(2) + '</h1>';
        } else if(line.startsWith('## ')) {
          headings.push({ level: 2, text: line.substring(3), anchor: 'heading-' + i });
          lines[i] = '<h2 id="heading-' + i + '">' + line.substring(3) + '</h2>';
        } else if(line.startsWith('### ')) {
          headings.push({ level: 3, text: line.substring(4), anchor: 'heading-' + i });
          lines[i] = '<h3 id="heading-' + i + '">' + line.substring(4) + '</h3>';
        }
      }
      
      // 生成目录HTML
      let tocHtml = '<h2>目录</h2><ul class="toc-list" style="list-style-type: none; padding: 0;">';
      for(const heading of headings) {
        const indent = heading.level > 1 ? 'style="margin-left:' + ((heading.level-1)*15) + 'px; font-size:' + (1.2 - (heading.level-1)*0.1) + 'em;"' : 'style="font-size:1.2em;"';
        tocHtml += '<li class="toc-list-item" ' + indent + '><a href="#' + heading.anchor + '" style="text-decoration: none; color: #007bff;">' + heading.text + '</a></li>';
      }
      tocHtml += '</ul>';
      
      return {
        tocHtml: tocHtml,
        modifiedMarkdown: lines.join('\n')
      };
    }
    
    // 获取Base64编码的内容并解码
    const encodedContent = "` + encodedContent + `";
    const markdownContent = decodeBase64UTF8(encodedContent);
    
    // 生成目录
    const { tocHtml, modifiedMarkdown } = generateTOC(markdownContent);
    
    // 设置目录
    document.getElementById('toc').innerHTML = tocHtml;
    
    // 使用修改后的Markdown渲染内容
    document.getElementById('content').innerHTML = marked.parse(modifiedMarkdown);

    // 确保JavaScript在DOM加载完成后执行
    document.addEventListener('DOMContentLoaded', function() {
      document.querySelectorAll('.toc-list-item').forEach(item => {
        item.addEventListener('click', function(e) {
          e.stopPropagation(); // 防止事件冒泡
          this.classList.toggle('open');
        });
      });
    });
  </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
