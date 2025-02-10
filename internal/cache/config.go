package cache

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/invopop/yaml"
	"io"
	"os"
)

type RedisConfig struct {
	Addr   string `json:"address" yaml:"address" env:"ADDRESS"`
	Passwd string `env-required:"true" json:"passwd" yaml:"passwd" env:"PASSWORD"`
	DB     int    `env-required:"true" json:"DBName" yaml:"DBName" env:"DB"`
}

func (rc *RedisConfig) loadConfigParam(filePath string) error {
	_, err := os.Stat(filePath)
	if !(err == nil || !os.IsNotExist(err)) {
		return fmt.Errorf("-> os.Stat: файл по указаному пути не найден %s", filePath)
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("-> os.OpenFile: ошибка при открытии файла %s: %w", filePath, err)
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("-> io.ReadAll: ошибка при чтении файла %s: %w", filePath, err)
	}

	err = yaml.Unmarshal(buf, rc)
	if err != nil {
		return fmt.Errorf("-> yaml.Unmarshal: ошибка при кодировании файла: %w", err)
	}

	err = cleanenv.UpdateEnv(rc)
	if err != nil {
		return fmt.Errorf("-> cleanenv.UpdateEnv: ошибка при обновлении параметроа из переменныз окружения%w", err)
	}

	return nil
}

func (rc *RedisConfig) UpdateEnvAddress() error {
	err := cleanenv.ReadEnv(rc)
	if err != nil {
		return fmt.Errorf("-> cleanenv.ReadEnv: ошибка загрузки параметров из переменных окружения: %w", err)
	}
	return nil
}
