package cli

import (
	"fmt"

	"github.com/eljamo/libpass/v5/asset"
	"github.com/eljamo/libpass/v5/config"
	"github.com/eljamo/libpass/v5/config/option"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const CustomConfigKeyPath = "custom_config_path"

func generateConfig(cmd *cobra.Command) (*config.Settings, error) {
	baseCfg, customCfg, err := loadJSONFiles(cmd)
	if err != nil {
		return nil, fmt.Errorf("loadJSONFiles error: %w", err)
	}

	flagCfg, err := getCmdFlagsAsJSON(cmd)
	if err != nil {
		return nil, fmt.Errorf("getCmdFlags error: %w", err)
	}

	return config.Generate(baseCfg, customCfg, flagCfg)
}

// loadJSONFiles loads the base config and the custom config from the JSON files
func loadJSONFiles(cmd *cobra.Command) (map[string]any, map[string]any, error) {
	customCfg, err := getCustomConfigJSON(cmd)
	if err != nil {
		return nil, nil, err
	}

	presetValue, err := getPresetValue(cmd, customCfg)
	if err != nil {
		return nil, nil, err
	}

	baseCfg, err := asset.GetJSONPreset(presetValue)
	if err != nil {
		return nil, nil, err
	}

	return baseCfg, customCfg, nil
}

func getCustomConfigJSON(cmd *cobra.Command) (map[string]any, error) {
	path, err := getCustomConfigPath(cmd)
	if err != nil {
		return nil, err
	}

	customCfgJSON, err := loadCustomConfig(path)
	if err != nil {
		return nil, err
	}

	return customCfgJSON, nil
}

func getCmdFlagsAsJSON(cmd *cobra.Command) (map[string]any, error) {
	flags := make(map[string]any)

	var err error
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if err != nil {
			return
		}

		var flagErr error
		switch flag.Value.Type() {
		case "string":
			flags[flag.Name], flagErr = cmd.Flags().GetString(flag.Name)
		case "int":
			flags[flag.Name], flagErr = cmd.Flags().GetInt(flag.Name)
		case "bool":
			flags[flag.Name], flagErr = cmd.Flags().GetBool(flag.Name)
		case "stringSlice":
			flags[flag.Name], flagErr = cmd.Flags().GetStringSlice(flag.Name)
		}

		if flagErr != nil {
			err = flagErr
		}
	})

	if err != nil {
		return nil, fmt.Errorf("getCmdFlagsAsJSON error: %w", err)
	}

	return flags, nil
}

func getCustomConfigPath(cmd *cobra.Command) (string, error) {
	return cmd.Flags().GetString(CustomConfigKeyPath)
}

func loadCustomConfig(path string) (map[string]any, error) {
	if path == "" {
		return nil, nil
	}

	return asset.LoadJSONFile(path)
}

func getPresetFromCustomConfig(customCfgJSON map[string]any) string {
	if customCfgJSON == nil {
		return ""
	}

	if preset, ok := customCfgJSON[option.PresetKey].(string); ok {
		return preset
	}

	return ""
}

func getPresetValue(cmd *cobra.Command, customJSONCfg map[string]any) (string, error) {
	var presetValue string
	presetFlag, presetArgPresent, err := checkPresetFlag(cmd)
	if err != nil {
		return "", err
	}

	presetFromCustomCfg := getPresetFromCustomConfig(customJSONCfg)

	if !presetArgPresent && presetFromCustomCfg != "" {
		presetValue = presetFromCustomCfg
	} else {
		presetValue = presetFlag
	}

	return presetValue, nil
}

func checkPresetFlag(cmd *cobra.Command) (string, bool, error) {
	presetArgPresent := isFlagSet(cmd, option.PresetKey)
	presetFlag, err := cmd.Flags().GetString(option.PresetKey)
	if presetArgPresent && (err != nil || presetFlag == "") {
		return "", false, fmt.Errorf("invalid %s flag: %w", option.PresetKey, err)
	}

	return presetFlag, presetArgPresent, nil
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
