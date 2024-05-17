package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/eljamo/libpass/v7/config/option"
	"github.com/eljamo/libpass/v7/service"
	"github.com/spf13/cobra"
)

var version = "1.13.0"

var rootCmd = &cobra.Command{
	Use:          "mempass",
	Version:      version,
	Short:        "A memorable password generator",
	Long:         `A memorable password generator, a CLI version of xkpasswd.net written in Go`,
	SilenceUsage: true, // do not print usage when an error occurs since usage is large and distracts from the error
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := generateConfig(cmd)
		if err != nil {
			return err
		}

		pgs, err := service.NewPasswordGeneratorService(cfg)
		if err != nil {
			return err
		}

		pws, err := pgs.Generate()
		if err != nil {
			return err
		}

		for _, p := range pws {
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
	ccss := strings.Join(option.Presets[:], ", ")
	pco := strings.Join(option.PaddingCharacterOptions[:], ", ")
	sco := strings.Join(option.SeparatorCharacterOptions[:], ", ")
	ptcss := strings.Join(option.PaddingTypes[:], ", ")
	sccss := strings.Join(option.DefaultSpecialCharacters[:], ", ")
	ttcss := strings.Join(option.TransformTypes[:], ", ")
	wlcss := strings.Join(option.WordLists, ", ")

	// Preset and Custom Config Flags
	rootCmd.Flags().String(
		CustomConfigPathKey, "",
		"custom config file path, you can use this to load a custom config. Such as ones generated by xkpasswd.net",
	)
	rootCmd.Flags().String(
		"preset", option.PresetDefault,
		fmt.Sprintf("use a built-in preset. Valid values: %s", ccss),
	)

	// Word List Flags
	rootCmd.Flags().String(
		"word_list", option.WordListEN,
		fmt.Sprintf("use a built-in list of words. Valid values: %s", wlcss),
	)

	// Passwords Flags
	rootCmd.Flags().Int(
		"num_passwords", 3,
		"number of passwords to generate, valid values: 1+",
	)

	// Word Flags
	rootCmd.Flags().Int(
		"num_words", 3,
		"number of words, valid values: 2+",
	)
	rootCmd.Flags().String(
		"case_transform", option.CaseTransformRandom,
		fmt.Sprintf("case transformation, allowed values: %s", ttcss),
	)
	rootCmd.Flags().Int(
		"word_length_min", 4,
		"minimum word length, valid values: 1+",
	)
	rootCmd.Flags().Int(
		"word_length_max", 8,
		"maximum word length, valid values: 1+",
	)

	// Separator Flags
	rootCmd.Flags().StringSlice(
		"separator_alphabet",
		[]string{},
		fmt.Sprintf("comma-separated list of characters to separate password parts, example values: %s", sccss),
	)
	rootCmd.Flags().String(
		"separator_character",
		option.SeparatorCharacterRandom,
		fmt.Sprintf("character to separate password parts, example values: %s", sco),
	)

	// Padding Flags
	rootCmd.Flags().Int(
		"pad_to_length", 0,
		"length to pad the password to, will be ignored if less than the generated password length, valid values: 0+",
	)
	rootCmd.Flags().String(
		"padding_character",
		option.PaddingCharacterRandom,
		fmt.Sprintf("character to pad the password with, example values: %s", pco),
	)
	rootCmd.Flags().Int(
		"padding_characters_after", 2,
		"number of characters to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().Int(
		"padding_characters_before", 2,
		"number of characters to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().Int(
		"padding_digits_after", 2,
		"number of digits to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().Int(
		"padding_digits_before", 2,
		"number of digits to pad before the password, valid values: 0+",
	)
	rootCmd.Flags().String(
		"padding_type", option.PaddingTypeFixed,
		fmt.Sprintf("padding type, allowed values: %s", ptcss),
	)
	rootCmd.Flags().StringSlice(
		"symbol_alphabet",
		[]string{},
		fmt.Sprintf("comma-separated list of characters to pad the password with, example values: %s", sccss),
	)
}
