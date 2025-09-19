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

// ANSI 颜色常量
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

// 执行 git 命令并返回结果
func runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 验证工单号格式 (必须 bcds-<数字> 或 bcds-<数字>-xxx)
func validateWorkItem(workItem string) bool {
	re := regexp.MustCompile(`^bcds-\d+(-[a-z0-9]+)*$`)
	return re.MatchString(workItem)
}

// CommitType 提交类型结构
type CommitType struct {
	code        string
	description string
}

// 获取提交类型列表（有序）
func getCommitTypes() []CommitType {
	return []CommitType{
		{"feat", "新功能 (A new feature)"},
		{"fix", "Bug修复 (A bug fix)"},
		{"docs", "文档更新 (Documentation only changes)"},
		{"style", "代码格式 (Changes that do not affect code meaning)"},
		{"refactor", "代码重构 (Neither fixes bug nor adds feature)"},
		{"perf", "性能优化 (A code change that improves performance)"},
		{"test", "测试相关 (Adding or correcting tests)"},
		{"build", "构建相关 (Affect build system or dependencies)"},
		{"ci", "CI配置 (Changes to CI configuration files)"},
		{"chore", "其他杂项 (Other changes)"},
	}
}

// 验证英文描述
func validateEnglishDescription(desc string) bool {
	if len(desc) == 0 || len(desc) > 72 {
		return false
	}
	// 首字母必须大写
	// if !strings.HasPrefix(strings.ToUpper(desc[:1]), desc[:1]) {
	// 	return false
	// }
	// 只允许英文和基本标点
	re := regexp.MustCompile(`^[A-Za-z0-9 ,.!?\-$begin:math:text$$end:math:text$]+$`)
	return re.MatchString(desc)
}

// 从控制台读取一行输入
func readLine(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// 从控制台读取多行输入，空行结束
func readMultiline(prompt string) string {
	fmt.Println(prompt + " (直接回车结束):")
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for {
		scanner.Scan()
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			break
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// 程序入口
func main() {
	// 检查是否有暂存的更改
	_, err := runGitCommand("diff", "--cached", "--quiet")
	if err == nil {
		fmt.Println(ColorRed + "错误：没有暂存的更改，请先使用 'git add' 添加文件。" + ColorReset)
		os.Exit(1)
	}

	fmt.Println(ColorCyan + "=== Git 提交助手 (Custom Commit CLI) ===" + ColorReset)

	// 1. 输入并验证工单号
	var workItem string
	for {
		workItem = strings.ToLower(readLine(ColorYellow + "1. 请输入工单号 (格式: bcds-<数字> 或 bcds-<数字>-xxx): " + ColorReset))
		if validateWorkItem(workItem) {
			break
		}
		fmt.Println(ColorRed + "工单号格式无效，请重新输入。" + ColorReset)
	}

	// 2. 选择提交类型
	commitTypes := getCommitTypes()
	fmt.Println("\n" + ColorYellow + "2. 请选择提交类型:" + ColorReset)
	for i, ct := range commitTypes {
		fmt.Printf("   %d. %-8s %s\n", i+1, ct.code, ct.description)
	}

	var commitType string
	for {
		input := readLine(fmt.Sprintf("   请输入编号 (1-%d): ", len(commitTypes)))
		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= len(commitTypes) {
			commitType = commitTypes[choice-1].code
			break
		}
		fmt.Printf(ColorRed+"请输入有效的数字 (1-%d)\n"+ColorReset, len(commitTypes))
	}

	// 3. 输入并验证英文简短描述
	fmt.Println("\n" + ColorYellow + "3. 请输入简短描述（英文，首字母大写，不超过72个字符）:" + ColorReset)
	var description string
	for {
		description = readLine("   描述: ")
		if validateEnglishDescription(description) {
			break
		}
		fmt.Println(ColorRed + "描述不符合规范，请确保使用英文、首字母大写且不超过72个字符。" + ColorReset)
	}

	// 4. 输入可选的详细描述
	body := readMultiline("\n" + ColorYellow + "4. 输入详细描述（可选，多行，支持中文）" + ColorReset)

	// 构建提交信息
	commitMessage := fmt.Sprintf("%s(%s): %s", commitType, workItem, description)
	if body != "" {
		commitMessage = commitMessage + "\n\n" + body
	}

	fmt.Println("\n" + ColorGreen + "生成的提交信息:" + ColorReset)
	fmt.Println("--------------------------------------------------")
	fmt.Println(commitMessage)
	fmt.Println("--------------------------------------------------")

	// 确认提交
	confirm := strings.ToLower(readLine("是否确认提交? (Y/N): "))
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println(ColorRed + "提交已取消。" + ColorReset)
		os.Exit(0)
	}

	// 执行提交
	output, err := runGitCommand("commit", "-m", commitMessage)
	if err != nil {
		fmt.Printf(ColorRed+"提交失败: %s\n"+ColorReset, output)
		os.Exit(1)
	}

	// 提交成功
	fmt.Println("\n" + ColorGreen + "=== 提交成功 ===" + ColorReset)
	fmt.Printf("提交内容:\n%s\n", commitMessage)
	fmt.Println("可以使用 'git push' 推送更改。")
}
