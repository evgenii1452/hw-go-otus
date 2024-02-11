package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

const testEnvDir = "./testdata/testenv"

func TestReadDir(t *testing.T) {
	createEnvs := func(vars map[string]string) error {
		err := exec.Command("mkdir", "-p", testEnvDir).Run()
		if err != nil {
			return fmt.Errorf("create env dir: %w", err)
		}

		for k, v := range vars {
			file, err := os.Create(fmt.Sprintf("%s/%s", testEnvDir, k))
			if err != nil {
				os.RemoveAll(testEnvDir)
				panic(fmt.Errorf("create env file: %w", err))
			}

			_, err = file.Write([]byte(v))
			if err != nil {
				os.RemoveAll(testEnvDir)
				panic(fmt.Errorf("write env file: %w", err))
			}
		}

		return nil
	}

	t.Run("test success", func(t *testing.T) {
		envMap := make(map[string]string)
		envMap["VAR1"] = "VALUE1"
		envMap["VAR2"] = "VALUE2"
		envMap["VAR3"] = ""

		err := createEnvs(envMap)
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(testEnvDir)

		envs, err := ReadDir(testEnvDir)
		require.Nil(t, err)

		for k, v := range envs {
			if k == "VAR3" {
				require.Equal(t, envMap[k], v.Value)
				require.True(t, v.NeedRemove)
				continue
			}

			require.Equal(t, envMap[k], v.Value)
			require.False(t, v.NeedRemove)
		}
	})

	t.Run("test success empty dir", func(t *testing.T) {
		err := createEnvs(make(map[string]string))
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(testEnvDir)

		envs, err := ReadDir(testEnvDir)
		require.Nil(t, err)
		require.Equal(t, 0, len(envs))
	})

	t.Run("test dir not exists", func(t *testing.T) {
		envs, err := ReadDir(testEnvDir)
		require.NotNil(t, err)
		require.Nil(t, envs)
	})
}
