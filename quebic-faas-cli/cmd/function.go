package cmd

import (
	"fmt"
	"os"
	"quebic-faas/quebic-faas-cli/common"
	"quebic-faas/types"
	"reflect"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var functionInputFile string
var functionStart bool

func init() {
	setupFunctionCmds()
	setupFunctionFlags()
}

var functionCmd = &cobra.Command{
	Use:   "function",
	Short: "Function commonds",
	Long:  `Function commonds`,
}

func setupFunctionCmds() {

	//function
	functionCmd.AddCommand(functionCreateCmd)
	functionCmd.AddCommand(functionUpdateCmd)
	functionCmd.AddCommand(functionGetALLCmd)
	functionCmd.AddCommand(functionInspectCmd)

	//function-container
	functionCmd.AddCommand(functionDeployCmd)
	functionCmd.AddCommand(functionStopCmd)

	//function-logs
	functionCmd.AddCommand(functionLogsCmd)

}

func setupFunctionFlags() {

	//function-create
	functionCreateCmd.PersistentFlags().StringVarP(&functionInputFile, "file", "f", "function.yml", "function input file")
	functionCreateCmd.PersistentFlags().BoolVarP(&functionStart, "start", "s", true, "if true function-container will start. otherwise not")

	//function-update
	functionUpdateCmd.PersistentFlags().StringVarP(&functionInputFile, "file", "f", "function.yml", "function input file")
	functionUpdateCmd.PersistentFlags().BoolVarP(&functionStart, "start", "s", true, "if true function-container will start. otherwise not")

	functionInspectCmd.PersistentFlags().StringVarP(&functionName, "name", "n", "", "function name")

}

var functionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "function : create",
	Long:  `function : create`,
	Run: func(cmd *cobra.Command, args []string) {
		functionSave(cmd, args, true)
	},
}

var functionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "function : update",
	Long:  `function : update`,
	Run: func(cmd *cobra.Command, args []string) {
		functionSave(cmd, args, false)
	},
}

var functionGetALLCmd = &cobra.Command{
	Use:   "ls",
	Short: "function : get-all",
	Long:  `function : get-all`,
	Run: func(cmd *cobra.Command, args []string) {
		functionGetALL(cmd, args)
	},
}

var functionInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "function : inspect function details",
	Long:  `function : inspect function details`,
	Run: func(cmd *cobra.Command, args []string) {
		functionGetByName(cmd, args)
	},
}

func functionSave(cmd *cobra.Command, args []string, isAdd bool) {

	functionDTO := &types.FunctionDTO{}
	err := common.ParseYAMLToObject(functionInputFile, functionDTO)
	if err != nil {
		prepareError(cmd, err)
	}

	mgrService := appContainer.GetMgrService()

	var errResponse *types.ErrorResponse
	if isAdd {
		errResponse = mgrService.FunctionCreate(functionDTO)
	} else {
		errResponse = mgrService.FunctionUpdate(functionDTO)
	}

	if errResponse != nil {
		prepareErrorResponse(cmd, errResponse)
	}

	color.Green("%s function : saved", functionDTO.Function.GetID())

	if functionStart {
		functionContainerOP(cmd, args, function_container_op_deploy, functionDTO.Function.GetID())
	}
}

func functionGetALL(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	functions, err := mgrService.FunctionsGetALL()
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Docker_Image_ID", "Route"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(prepareFunctionsTable(functions))
	table.Render()

}

func functionGetByName(cmd *cobra.Command, args []string) {

	mgrService := appContainer.GetMgrService()
	function, err := mgrService.FunctionsGetByName(functionName)
	if err != nil {
		prepareErrorResponse(cmd, err)
	}

	ymlStr, _ := yaml.Marshal(function)
	fmt.Printf("%s", ymlStr)

}

func prepareFunctionsTable(data []types.Function) [][]string {

	var rows [][]string

	for _, val := range data {

		name := val.Name
		dockerImageID := val.DockerImageID
		route := val.Route

		rows = append(rows, []string{name, dockerImageID, route})

	}

	return rows

}

func getField(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}