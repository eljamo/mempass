package cli

import (
	"encoding/json"
	"fmt"

	"github.com/eljamo/mempass/asset"
	"github.com/eljamo/mempass/internal/config"
	"github.com/eljamo/mempass/internal/json_merge"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	CUSTOM_CONFIG_PATH_KEY string = "custom_config_path"
)

func mergeConfigs(cmd *cobra.Command) (string, error) {
	baseCfg, err := getBaseJSONPreset(cmd)
	if err != nil {
		return "", err
	}

	flagData := make(map[string]interface{})
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		switch flag.Value.Type() {
		case "string":
			flagData[flag.Name], _ = cmd.Flags().GetString(flag.Name)
		case "int":
			flagData[flag.Name], _ = cmd.Flags().GetInt(flag.Name)
		case "bool":
			flagData[flag.Name], _ = cmd.Flags().GetBool(flag.Name)
		case "stringSlice":
			flagData[flag.Name], _ = cmd.Flags().GetStringSlice(flag.Name)
		}
	})

	flagCfg, err := json.Marshal(flagData)
	if err != nil {
		return "", err
	}

	return json_merge.Merge(baseCfg, string(flagCfg))
}

func getBaseJSONPreset(cmd *cobra.Command) (string, error) {
	cfgName, err := cmd.Flags().GetString(config.PRESET_KEY)
	if err != nil {
		return "", fmt.Errorf("failed to get %s flag: %w", config.PRESET_KEY, err)
	}

	cfgJSON, err := asset.GetJSONPreset(cfgName)
	if err != nil {
		return "", fmt.Errorf("failed to get JSON config: %w", err)
	}

	path, err := cmd.Flags().GetString(CUSTOM_CONFIG_PATH_KEY)
	if err != nil {
		return "", fmt.Errorf("failed to get %s flag: %w", CUSTOM_CONFIG_PATH_KEY, err)
	}

	if path != "" {
		customCfgJSON, err := asset.LoadJSONFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to load JSON file from path %s: %w", path, err)
		}
		return json_merge.Merge(cfgJSON, customCfgJSON)
	}

	return cfgJSON, nil
}

func generateConfig(cmd *cobra.Command, args []string) (*config.Config, error) {
	var cfg *config.Config
	cfgJSON, err := mergeConfigs(cmd)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal([]byte(cfgJSON), &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
