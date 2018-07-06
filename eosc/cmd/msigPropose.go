// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/msig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// msigProposeCmd represents the msigPropose command
var msigProposeCmd = &cobra.Command{
	Use:   "propose [proposer] [proposal name] [transaction_file.json]",
	Short: "Propose a new transaction in the eosio.msig contract",
	Long: `Propose a new transaction in the eosio.msig contract

Pass --requested-permissions
`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()

		proposer := toAccount(args[0], "proposer")
		proposalName := toName(args[1], "proposal name")
		transactionFileName := args[2]

		cnt, err := ioutil.ReadFile(transactionFileName)
		errorCheck("reading transaction file", err)

		var tx *eos.Transaction
		err = json.Unmarshal(cnt, &tx)
		errorCheck("parsing transaction file", err)

		requested, err := permissionsToPermissionLevels(viper.GetStringSlice("msig-propose-cmd-requested-permissions"))
		errorCheck("requested permissions", err)
		if len(requested) == 0 {
			errorCheck("--requested-permissions", errors.New("missing values"))
		}

		pushEOSCActions(api,
			msig.NewPropose(proposer, proposalName, requested, tx),
		)
	},
}

func init() {
	msigCmd.AddCommand(msigProposeCmd)

	msigProposeCmd.Flags().StringSliceP("requested-permissions", "", []string{}, "Permissions requested, specify multiple times or separated by a comma.")

	for _, flag := range []string{"requested-permissions"} {
		if err := viper.BindPFlag("msig-propose-cmd-"+flag, msigProposeCmd.Flags().Lookup(flag)); err != nil {
			panic(err)
		}
	}

}
