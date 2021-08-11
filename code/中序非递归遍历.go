package code

import (
	"fmt"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func inorderTraversal(root *TreeNode) []int {
	var res []int
	stack := make([]*TreeNode, 0)
	for root != nil || len(stack) > 0 {
		for root != nil {
			stack = append(stack, root)
			root = root.Left
		}
		root = stack[len(stack)-1]
		fmt.Print(&root)
		fmt.Print(&stack[len(stack)-1])
		stack = stack[:len(stack)-1]
		res = append(res, root.Val)
		root = root.Right

		// res = append(res, stack[len(stack) - 1].Val)
		// stack = stack[:len(stack) - 1]
		// root = root.Right
	}
	return res
}
