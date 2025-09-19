package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// 执行git命令并返回结果
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 验证工单号格式
func validateWorkItem(workItem string) bool {
	// 可根据团队规范修改正则表达式
	re := regexp.MustCompile(`^([A-Za-z]+-)?\d+$`)
	return re.MatchString(workItem)
}

// 获取提交类型列表
func getCommitTypes() map[int]struct {
	code        string
	description string
} {
	return map[int]struct {
		code        string
		description string
	}{
		1:  {code: "feat", description: "新功能 (A new feature)"},
		2:  {code: "fix", description: "Bug修复 (A bug fix)"},
		3:  {code: "docs", description: "文档更新 (Documentation only changes)"},
		4:  {code: "style", description: "代码格式 (Changes that do not affect code meaning)"},
		5:  {code: "refactor", description: "代码重构 (Neither fixes bug nor adds feature)"},
		6:  {code: "perf", description: "性能优化 (A code change that improves performance)"},
		7:  {code: "test", description: "测试相关 (Adding or correcting tests)"},
		8:  {code: "build", description: "构建相关 (Affect build system or dependencies)"},
		9:  {code: "ci", description: "CI配置 (Changes to CI configuration files)"},
		10: {code: "chore", description: "其他杂项 (Other changes)"},
	}
}

// 验证英文描述
func validateEnglishDescription(desc string) bool {
	if len(desc) == 0 {
		return false
	}
	if len(desc) > 72 {
		return false
	}
	// 首字母必须大写
	if len(desc) > 0 && !strings.HasPrefix(strings.ToUpper(desc[:1]), desc[:1]) {
		return false
	}
	// 只允许英文和基本标点
	re := regexp.MustCompile(`^[A-Za-z0-9 ,.!?\-\(\)]+$`)
	return re.MatchString(desc)
}

// 读取用户输入
func readInput(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
	// 检查是否有暂存的更改
	_, err := runGitCommand("diff", "--cached", "--quiet")
	if err == nil {
		fmt.Println("错误：没有暂存的更改，请先使用 'git add' 添加文件。")
		os.Exit(1)
	}

	fmt.Println("=== 自定义提交流程 ===")

	// 1. 输入并验证工单号
	var workItem string
	for {
		workItem = readInput("1. 请输入工单号 (例如 PROJ-1234): ")
		if validateWorkItem(workItem) {
			break
		}
		fmt.Println("工单号格式无效，请重新输入。")
	}

	// 2. 选择提交类型
	commitTypes := getCommitTypes()
	fmt.Println("\n2. 请选择提交类型:")
	for i, ct := range commitTypes {
		fmt.Printf("   %d. %-8s %s\n", i, ct.code, ct.description)
	}

	var commitType string
	for {
		input := readInput(fmt.Sprintf("   请输入编号 (1-%d): ", len(commitTypes)))
		choice, err := strconv.Atoi(input)
		if err == nil {
			if ct, exists := commitTypes[choice]; exists {
				commitType = ct.code
				break
			}
		}
		fmt.Printf("请输入有效的数字 (1-%d)\n", len(commitTypes))
	}

	// 3. 输入并验证英文描述
	fmt.Println("\n3. 请输入修改内容（英文，首字母大写，不超过72个字符）:")
	var description string
	for {
		description = readInput("   描述: ")
		if validateEnglishDescription(description) {
			break
		}
		fmt.Println("描述不符合规范，请确保使用英文、首字母大写且不超过72个字符。")
	}

	// 构建提交信息
	commitMessage := fmt.Sprintf("%s(%s): %s", commitType, workItem, description)
	fmt.Printf("\n生成的提交信息:\n%s\n", commitMessage)

	// 确认提交
	confirm := readInput("是否确认提交? (y/n): ")
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("提交已取消。")
		os.Exit(0)
	}

	// 执行提交
	output, err := runGitCommand("commit", "-m", commitMessage)
	if err != nil {
		fmt.Printf("提交失败: %s\n", output)
		os.Exit(1)
	}

	// 提交成功后显示最终的commit信息
	fmt.Println("\n=== 提交成功 ===")
	fmt.Printf("提交内容: %s\n", commitMessage)
	fmt.Println("可以使用 'git push' 推送更改。")
}
