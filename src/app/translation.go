package app

import (
	"fmt"
	"path/filepath"
)

var TranslationDefaultLanguageKey = "TRANSLATION_DEFAULT_LANGUAGE"

type (
	Translation struct {
		FilePath            string   `mapstructure:"TRANSLATION_FILE_PATH"`
		LanguagePreferences []string `mapstructure:"TRANSLATION_LANG_PREFERENCES"`
		DefaultLanguage     string   `mapstructure:"TRANSLATION_DEAULT_LANG"`
	}
)

func (cfg Translation) TranslationJSONFiles() []string {
	var files []string
	languages := append(cfg.LanguagePreferences, cfg.DefaultLanguage)
	for _, lang := range languages {
		fileName := fmt.Sprintf("%s.all.json", lang)
		files = append(files, filepath.Join(cfg.FilePath, fileName))
	}
	return files
}
