package http

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/invopop/yaml"
	"io"
	"os"
)

type ServerAddress struct {
	Localhost   string `json:"localhost" yaml:"localhost"`
	DefaultPort int    `json:"defaultPort" yaml:"defaultPort"`
	EnvAddress  string `env-required:"true" json:"envAddress" yaml:"envAddress" env:"SERVER_ADDRESS"`
}

func (a *ServerAddress) LoadConfigAddress(filePath string) error {
	_, err := os.Stat(filePath)
	if !(err == nil || !os.IsNotExist(err)) {
		return fmt.Errorf("-> os.Stat: файла не существует %s: %w", filePath, err)
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("-> os.OpenFile: ошибка при открытии файла %s: %w", filePath, err)

	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("-> os.OpenFile: ошибка при чтении файла %s: %w", filePath, err)
	}

	err = yaml.Unmarshal(buf, a)
	if err != nil {
		return fmt.Errorf("->  yaml.Unmarshal: ошибка при конвертации: %w", err)
	}

	return nil
}

func (a *ServerAddress) UpdateEnvAddress() error {
	err := cleanenv.ReadEnv(a)
	if err != nil {
		return fmt.Errorf("-> cleanenv.ReadEnv: ошибка загрузки параметров из переменных окружения: %w", err)
	}
	return nil
}
