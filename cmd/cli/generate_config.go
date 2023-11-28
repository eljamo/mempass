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

const CustomConfigKeyPath = "custom_config_path"

func mergeConfigs(cmd *cobra.Command) (string, error) {
	baseCfg, err := getBaseJSONPreset(cmd)
	if err != nil {
		return "", fmt.Errorf("getBaseJSONPreset error: %w", err)
	}

	flagData, err := getFlagData(cmd)
	if err != nil {
		return "", fmt.Errorf("getFlagData error: %w", err)
	}

	flagCfg, err := json.Marshal(flagData)
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}

	return json_merge.Merge(baseCfg, string(flagCfg))
}

func getFlagData(cmd *cobra.Command) (map[string]any, error) {
	flagData := make(map[string]any)
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		var err error
		switch flag.Value.Type() {
		case "string":
			flagData[flag.Name], err = cmd.Flags().GetString(flag.Name)
		case "int":
			flagData[flag.Name], err = cmd.Flags().GetInt(flag.Name)
		case "bool":
			flagData[flag.Name], err = cmd.Flags().GetBool(flag.Name)
		case "stringSlice":
			flagData[flag.Name], err = cmd.Flags().GetStringSlice(flag.Name)
		}
		if err != nil {
			return
		}
	})

	return flagData, nil
}

func isFlagSet(cmd *cobra.Command, flagKey string) bool {
	var flagSet bool
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Name == flagKey {
			flagSet = true
		}
	})

	return flagSet
}

func checkPresetFlag(cmd *cobra.Command) (string, bool, error) {
	presetArgPresent := isFlagSet(cmd, config.PresetKey)
	presetFlag, err := cmd.Flags().GetString(config.PresetKey)
	if presetArgPresent && (err != nil || presetFlag == "") {
		return "", false, fmt.Errorf("invalid %s flag: %w", config.PresetKey, err)
	}

	return presetFlag, presetArgPresent, nil
}

func getCustomConfigPath(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(CustomConfigKeyPath)
}

func loadCustomConfig(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	return asset.LoadJSONFile(path)
}

func getPresetFromCustomConfig(customCfgJSON string) (string, error) {
	if customCfgJSON == "" {
		return "", nil
	}

	var customCfg map[string]any
	if err := json.Unmarshal([]byte(customCfgJSON), &customCfg); err != nil {
		return "", fmt.Errorf("json unmarshal error for custom config: %w", err)
	}

	if preset, ok := customCfg[config.PresetKey].(string); ok {
		return preset, nil
	}

	return "", nil
}

func getBaseJSONPreset(cmd *cobra.Command) (string, error) {
	presetFlag, presetArgPresent, err := checkPresetFlag(cmd)
	if err != nil {
		return "", err
	}

	path, err := getCustomConfigPath(cmd)
	if err != nil {
		return "", err
	}

	customCfgJSON, err := loadCustomConfig(path)
	if err != nil {
		return "", err
	}

	presetFromCustomCfg, err := getPresetFromCustomConfig(customCfgJSON)
	if err != nil {
		return "", err
	}

	cfgName := presetFlag
	if !presetArgPresent && presetFromCustomCfg != "" {
		cfgName = presetFromCustomCfg
	}

	cfgJSON, err := asset.GetJSONPreset(cfgName)
	if err != nil {
		return "", fmt.Errorf("failed to get JSON config: %w", err)
	}

	if path != "" && customCfgJSON != "" {
		return json_merge.Merge(cfgJSON, customCfgJSON)
	}

	return cfgJSON, nil
}

func mergeCustomConfig(cfgJSON, path string) (string, error) {
	customCfgJSON, err := asset.LoadJSONFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load JSON file from path %s: %w", path, err)
	}

	return json_merge.Merge(cfgJSON, customCfgJSON)
}

func generateConfig(cmd *cobra.Command, args []string) (*config.Config, error) {
	cfgJSON, err := mergeConfigs(cmd)
	if err != nil {
		return nil, fmt.Errorf("mergeConfigs error: %w", err)
	}

	var cfg config.Config
	if err := json.Unmarshal([]byte(cfgJSON), &cfg); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &cfg, nil
}
