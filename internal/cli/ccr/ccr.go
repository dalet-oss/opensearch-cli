/*
Copyright © 2025 Sergei Iakovlev syakovlev@dalet.com
*/
package ccr

import "github.com/spf13/cobra"

func NewCCRCmd() *cobra.Command {
	//
	return ccrCmd
}

var ccrCmd = &cobra.Command{}
