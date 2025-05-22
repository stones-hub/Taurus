package main

import (
	"Taurus/pkg/validate"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// User 用户数据结构
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,gte=18,lte=120"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin user guest"`
	Mobile   string `json:"mobile" validate:"omitempty,len=11"`
}

// Product 产品数据结构
type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	Status      string  `json:"status" validate:"required,oneof=active inactive deleted"`
}

// simpleValidateExample 简单验证示例
func simpleValidateExample() {
	fmt.Println("\n=== 结构体验证 ===")
	validateStructExample()

	fmt.Println("\n=== 单个变量验证 ===")
	validateSingleVarExample()

	fmt.Println("\n=== 映射验证 ===")
	validateMapExample()

	fmt.Println("\n=== 切片验证 ===")
	validateSliceExample()

	fmt.Println("\n=== 自定义验证规则 ===")
	customValidationExample()
}

// validateStructExample 结构体验证示例
func validateStructExample() {
	// 创建一个有效的用户
	validUser := &User{
		Username: "zhangsan",
		Email:    "zhangsan@example.com",
		Age:      30,
		Password: "password123",
		Role:     "admin",
	}

	// 验证有效用户
	if err := validate.ValidateStruct(validUser); err != nil {
		fmt.Printf("有效用户验证失败: %v\n", err)
	} else {
		fmt.Println("有效用户验证通过")
	}

	// 创建一个无效的用户
	invalidUser := &User{
		Username: "li",            // 太短
		Email:    "invalid-email", // 无效的邮箱
		Age:      15,              // 年龄太小
		Password: "123",           // 密码太短
		Role:     "superuser",     // 不在允许的角色列表中
		Mobile:   "1234",          // 长度不符合要求
	}

	// 验证无效用户
	if err := validate.ValidateStruct(invalidUser); err != nil {
		fmt.Printf("无效用户验证失败: %v\n", err)

		// 获取验证错误详情
		if valErrs, ok := err.(validate.ValidationErrors); ok {
			fmt.Println("\n详细错误信息:")
			for _, e := range valErrs {
				fmt.Printf("字段: %s, 标签: %s, 值: %v, 错误: %s\n",
					e.Field, e.Tag, e.Value, e.Message)
			}

			// 获取字段错误映射
			fieldErrors := validate.GetFieldErrors(valErrs)
			fmt.Println("\n字段错误映射:")
			for field, msg := range fieldErrors {
				fmt.Printf("%s: %s\n", field, msg)
			}
		}
	}
}

// validateSingleVarExample 单个变量验证示例
func validateSingleVarExample() {
	// 验证电子邮件
	email := "invalid-email"
	if err := validate.ValidateVar(email, "email"); err != nil {
		fmt.Printf("无效的电子邮件: %v\n", err)
	}

	// 验证有效的电子邮件
	validEmail := "test@example.com"
	if err := validate.ValidateVar(validEmail, "email"); err != nil {
		fmt.Printf("电子邮件验证失败: %v\n", err)
	} else {
		fmt.Printf("有效的电子邮件: %s\n", validEmail)
	}

	// 验证长度
	password := "123"
	if err := validate.ValidateVar(password, "min=8"); err != nil {
		fmt.Printf("密码太短: %v\n", err)
	}

	// 验证数字范围
	age := 15
	if err := validate.ValidateVar(age, "gte=18,lte=120"); err != nil {
		fmt.Printf("年龄不在有效范围内: %v\n", err)
	}

	// 验证变量之间的关系
	min := 10
	max := 5
	if err := validate.ValidateVarWithValue(max, min, "gtefield"); err != nil {
		fmt.Printf("最大值必须大于等于最小值: %v\n", err)
	}
}

// validateMapExample 映射验证示例
func validateMapExample() {
	// 创建包含用户的映射
	userMap := map[string]interface{}{
		"user1": &User{
			Username: "zhangsan",
			Email:    "zhangsan@example.com",
			Age:      30,
			Password: "password123",
			Role:     "admin",
		},
		"user2": &User{
			Username: "li", // 太短
			Email:    "invalid-email",
			Age:      15,
			Password: "123",
			Role:     "superuser",
		},
	}

	if err := validate.ValidateMap(userMap); err != nil {
		fmt.Printf("映射验证失败: %v\n", err)
		if valErrs, ok := err.(validate.ValidationErrors); ok {
			for _, e := range valErrs {
				fmt.Printf("字段: %s, 错误: %s\n", e.Field, e.Message)
			}
		}
	}

	// 创建包含产品的映射
	productMap := map[int]Product{
		1: {
			Name:       "手机",
			Price:      1999.99,
			CategoryID: 1,
			Status:     "active",
		},
		2: {
			Name:       "A",       // 名称太短
			Price:      -1,        // 无效的价格
			CategoryID: 0,         // 无效的类别ID
			Status:     "pending", // 无效的状态
		},
	}

	if err := validate.ValidateMap(productMap); err != nil {
		fmt.Printf("产品映射验证失败: %v\n", err)
	}
}

// validateSliceExample 切片验证示例
func validateSliceExample() {
	// 创建用户切片
	users := []*User{
		{
			Username: "zhangsan",
			Email:    "zhangsan@example.com",
			Age:      30,
			Password: "password123",
			Role:     "admin",
		},
		{
			Username: "li", // 太短
			Email:    "invalid-email",
			Age:      15,
			Password: "123",
			Role:     "superuser",
		},
	}

	if err := validate.ValidateSlice(users); err != nil {
		fmt.Printf("切片验证失败: %v\n", err)
		if valErrs, ok := err.(validate.ValidationErrors); ok {
			for _, e := range valErrs {
				fmt.Printf("字段: %s, 错误: %s\n", e.Field, e.Message)
			}
		}
	}

	// 创建产品切片
	products := []Product{
		{
			Name:       "手机",
			Price:      1999.99,
			CategoryID: 1,
			Status:     "active",
		},
		{
			Name:       "A",       // 名称太短
			Price:      -1,        // 无效的价格
			CategoryID: 0,         // 无效的类别ID
			Status:     "pending", // 无效的状态
		},
	}

	if err := validate.ValidateSlice(products); err != nil {
		fmt.Printf("产品切片验证失败: %v\n", err)
	}
}

// customValidationExample 自定义验证规则示例
func customValidationExample() {
	// 创建验证器
	v := validate.New()

	// 注册自定义验证规则：中国手机号
	if err := v.RegisterCustomValidation("validate_mobile", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		// 简单的中国手机号验证：11位数字，以1开头
		if len(value) != 11 || value[0] != '1' {
			return false
		}
		for _, c := range value {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}, "{0}必须是有效的中国手机号码"); err != nil {
		fmt.Printf("注册自定义验证规则失败: %v\n", err)
		return
	}

	// 定义使用自定义验证规则的结构体
	type Contact struct {
		Mobile string `json:"mobile" validate:"required,validate_mobile"` // 使用自定义验证规则validate_mobile
		Tel    string `json:"tel" validate:"omitempty"`
	}

	// 验证有效的手机号
	validContact := Contact{Mobile: "13812345678"}
	errors, _ := v.ValidateStruct(validContact)
	if len(errors) > 0 {
		fmt.Printf("有效联系人验证失败: %v\n", errors)
	} else {
		fmt.Println("有效联系人验证通过")
	}

	// 验证无效的手机号
	invalidContact := &Contact{Mobile: "123456"}
	errors, _ = v.ValidateStruct(invalidContact)
	if len(errors) > 0 {
		fmt.Printf("无效联系人验证失败: %v\n", errors)
		for _, e := range errors {
			fmt.Printf("字段: %s, 标签: %s, 值: %v, 错误: %s\n",
				e.Field, e.Tag, e.Value, e.Message)
		}
	}
}

func main() {
	fmt.Println("=== 验证示例开始 ===")
	simpleValidateExample()
	fmt.Println("=== 验证示例结束 ===")
}
