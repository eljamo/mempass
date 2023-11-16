package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/eljamo/mempass/internal/config"
	"github.com/eljamo/mempass/internal/service"
	"github.com/spf13/cobra"
)

var version = "1.0.2"

var rootCmd = &cobra.Command{
	Use:     "mempass",
	Version: version,
	Short:   "A memorable password generator",
	Long:    `A memorable password generator, a CLI version of xkpasswd.net written in Go`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := generateConfig(cmd, args)
		if err != nil {
			return err
		}

		ps := service.NewPasswordGeneratorService(cfg)
		pw, err := ps.Generate()
		if err != nil {
			return err
		}

		for _, p := range pw {
			fmt.Println(p)
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	ccss := strings.Join(config.Preset[:], ", ")
	pcasccss := strings.Join(config.PaddingCharacterAndSeparatorCharacter[:], ", ")
	ptcss := strings.Join(config.PaddingType[:], ", ")
	sccss := strings.Join(config.DefaultSpecialCharacters[:], ", ")
	ttcss := strings.Join(config.TransformType[:], ", ")
	wlcss := strings.Join(config.WordLists, ", ")

	rootCmd.Flags().String(
		"preset", config.DEFAULT,
		fmt.Sprintf("use a built-in preset. Valid values: %s", ccss),
	)
	rootCmd.Flags().String(
		"custom_config_path", "",
		"custom config file path, you can use this to load a custom config. Such as ones generated by xkpasswd.net",
	)
	rootCmd.Flags().String(
		"word_list", config.EN,
		fmt.Sprintf("use a built-in list of words. Valid values: %s", wlcss),
	)

	rootCmd.Flags().Int(
		"num_passwords", 3,
		"number of passwords to generate, valid values: 1+",
	)
	rootCmd.Flags().Int(
		"num_words", 3,
		"number of words, valid values: 2+",
	)
	rootCmd.Flags().Int(
		"word_length_min", 4,
		"minimum word length, valid values: 1+",
	)
	rootCmd.Flags().Int(
		"word_length_max", 8,
		"maximum word length, valid values: 1+",
	)
	rootCmd.Flags().String(
		"case_transform", config.RANDOM,
		fmt.Sprintf("case transformation, allowed values: %s", ttcss),
	)
	rootCmd.Flags().String(
		"separator_character",
		config.RANDOM,
		fmt.Sprintf("character to separate password parts, example values: %s", pcasccss),
	)
	rootCmd.Flags().StringSlice(
		"separator_alphabet",
		[]string{},
		fmt.Sprintf("comma-separated list of characters to separate password parts, example values: %s", sccss),
	)
	rootCmd.Flags().Int(
		"padding_digits_before", 2,
		"number of digits to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().Int(
		"padding_digits_after", 2,
		"number of digits to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().String(
		"padding_type", "",
		fmt.Sprintf("padding type, allowed values: %s", ptcss),
	)
	rootCmd.Flags().String(
		"padding_character",
		config.RANDOM,
		fmt.Sprintf("character to pad the password with, example values: %s", pcasccss),
	)
	rootCmd.Flags().StringSlice(
		"symbol_alphabet",
		[]string{},
		fmt.Sprintf("comma-separated list of characters to pad the password with, example values: %s", sccss),
	)
	rootCmd.Flags().Int(
		"pad_to_length", 0,
		"length to pad the password to, will be ignored if less than the generated password length, valid values: 0+",
	)
	rootCmd.Flags().Int(
		"padding_characters_before", 2,
		"number of characters to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().Int(
		"padding_characters_after", 2,
		"number of characters to pad before the password, valid values: 0+",
	)
}
