package cli

import (
	"fmt"

	"github.com/eljamo/libpass/v7/asset"
	"github.com/eljamo/libpass/v7/config"
	"github.com/eljamo/libpass/v7/config/option"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const CustomConfigPathKey string = "custom_config_path"

func generateConfig(cmd *cobra.Command) (*config.Settings, error) {
	customCfg, err := loadCustomConfig(cmd)
	if err != nil {
		return nil, err
	}

	flagCfg, err := getCmdFlags(cmd)
	if err != nil {
		return nil, err
	}

	presetValue, err := getPresetValue(cmd, customCfg)
	if err != nil {
		return nil, err
	}

	if presetValue == option.PresetDefault {
		return config.New(customCfg, flagCfg)
	}

	basePreset, err := loadBasePreset(presetValue)
	if err != nil {
		return nil, err
	}

	return config.New(basePreset, customCfg, flagCfg)
}

// Loads the base preset and the custom config from the JSON files
func loadCustomConfig(cmd *cobra.Command) (map[string]any, error) {
	customCfg, err := getCustomConfigJSON(cmd)
	if err != nil {
		return nil, err
	}

	return customCfg, nil
}

// Loads the base preset and the custom config from the JSON files
func loadBasePreset(presetValue string) (map[string]any, error) {
	basePreset, err := asset.GetJSONPreset(presetValue)
	if err != nil {
		return nil, err
	}

	return basePreset, nil
}

// Loads the custom config JSON file
func getCustomConfigJSON(cmd *cobra.Command) (map[string]any, error) {
	path, err := cmd.Flags().GetString(CustomConfigPathKey)
	if err != nil {
		return nil, err
	}

	customCfgJSON, err := loadCustomConfigJSON(path)
	if err != nil {
		return nil, err
	}

	return customCfgJSON, nil
}

// Returns a map of the cmd flags and their values
func getCmdFlags(cmd *cobra.Command) (map[string]any, error) {
	flags := make(map[string]any)

	var err error
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if err != nil {
			return
		}

		switch flag.Value.Type() {
		case "string":
			flags[flag.Name], err = cmd.Flags().GetString(flag.Name)
		case "int":
			flags[flag.Name], err = cmd.Flags().GetInt(flag.Name)
		case "bool":
			flags[flag.Name], err = cmd.Flags().GetBool(flag.Name)
		case "stringSlice":
			flags[flag.Name], err = cmd.Flags().GetStringSlice(flag.Name)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing cmd flags (%w)", err)
	}

	return flags, nil
}

// Loads the custom config JSON file
func loadCustomConfigJSON(path string) (map[string]any, error) {
	if path == "" {
		return nil, nil
	}

	return asset.LoadJSONFile(path)
}

// Returns the preset value from the custom config if it exists
func getPresetFromCustomConfig(customCfgJSON map[string]any) string {
	if customCfgJSON == nil {
		return ""
	}

	if preset, ok := customCfgJSON[option.ConfigKeyPreset].(string); ok {
		return preset
	}

	return ""
}

// Returns the preset value
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

// Returns the preset flag value and if preset flag was explicitly set
func checkPresetFlag(cmd *cobra.Command) (string, bool, error) {
	presetArgPresent := isFlagSet(cmd, option.ConfigKeyPreset)
	presetFlagValue, err := cmd.Flags().GetString(option.ConfigKeyPreset)
	if presetArgPresent && (err != nil || presetFlagValue == "") {
		return "", false, fmt.Errorf("invalid %s flag (%w)", option.ConfigKeyPreset, err)
	}

	return presetFlagValue, presetArgPresent, nil
}

// Checks if a Cobra flag has been explicitly set
func isFlagSet(cmd *cobra.Command, flagKey string) bool {
	var flagSet bool
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Name == flagKey {
			flagSet = true
		}
	})

	return flagSet
}
