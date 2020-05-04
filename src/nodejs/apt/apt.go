package apt
import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os"
	
	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface{
	LinkDirectoryInDepDir(string, string) error
	WriteProfileD(string, string) error
	CacheDir() String
	DepDir() String
}

type Command interface {
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(string, string, ...string) (string, error)
}

type Supplier struct{
	Stager         Stager
	Command        Command
}

func AptSetup(installDir String) (error,[]String,[]String,String)
{
	cacheDir:=s.Stager.CacheDir()
	aptCacheDir:=filepath.Join(cacheDir, "apt", "cache")
	stateDir := filepath.Join(cacheDir, "apt", "state")
	preferences := filepath.Join(cacheDir, "apt", "etc", "preferences")
	archiveDir :=filepath.Join(aptCacheDir, "archives")
	rootDir := "/etc/apt"

	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return err,nil,nil,nil
	}
	if err := os.MkdirAll(aptCacheDir, os.ModePerm); err != nil {
		return err,nil,nil,nil
	}
	if err := os.MkdirAll(stateDir, os.ModePerm); err != nil {
		return err,nil,nil,nil
	}
	if err := os.MkdirAll(installDir, os.ModePerm); err != nil {
		return err,nil,nil,nil
	}

	aptPrefs := filepath.Join(rootDir, "preferences")
	if exists, err := libbuildpack.FileExists(aptPrefs); err != nil {
		return err,nil,nil,nil
	} else if exists {
		if err := libbuildpack.CopyFile(aptPrefs, preferences); err != nil {
			return err,nil,nil,nil
		}
	} else {
		dirPath := filepath.Dir(preferences)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err,nil,nil,nil
		}
	}
	if err := os.MkdirAll(archiveDir, os.ModePerm); err != nil {
		return err,nil,nil,nil
	}

	sourcelist:=filepath.Join(cacheDir,"apt","sources","sources.list")
	aptSources:= filepath.Join(rootDir, "sources.list")
	if err := libbuildpack.CopyFile(aptSources, sourcelist); err != nil {
		return err,nil,nil,nil
	}
	options:=[]string{
			"-o", "debug::nolocking=true",
			"-o", "dir::cache=" + aptCacheDir,
			"-o", "dir::etc::sourcelist=" + sourcelist,
			"-o", "dir::state=" + stateDir,
			"-o", "Dir::Etc::preferences=" + preferences}
	doptions := append(options, "-y", "--allow-downgrades", "--allow-remove-essential", "--allow-change-held-packages", "-d", "install", "--reinstall")
	return (nil,options,doptions,archiveDir)
}

func AptUpdate(options String[],command Command)
{
	uargs := append(options, "update")	
	var errBuff bytes.Buffer
	if err := command.Execute("/", &errBuff, &errBuff, "apt-get", uargs...); err != nil {
		return fmt.Errorf("failed to apt-get update %s\n\n%s", errBuff.String(), err)
	}
	return nil;
}

func InstallPackages(archiveDir String,installDir String,command libbuildpack.command) error
{
	files, err := filepath.Glob(filepath.Join(archiveDir, "*.deb"))
	if err != nil {
		return err
	}
	for _, file := range files {
		output, err := command.Output("/", "dpkg", "-x", file, installDir)
	if err != nil {
		return fmt.Errorf("failed to install pkg %s\n\n%s\n\n%s", file, output, err.Error())
	}
	}
	return nil;
}

func LinkPackages(installDir String,stager supply.Stager) error
{
	for _, dirs := range [][]string{
		{"usr/bin", "bin"},
		{"usr/lib", "lib"},
		{"usr/lib/i386-linux-gnu", "lib"},
		{"usr/lib/x86_64-linux-gnu", "lib"},
		{"lib/x86_64-linux-gnu", "lib"},
		{"usr/include", "include"},
	} {
		dest := filepath.Join(installDir, dirs[0])
		if exists, err := libbuildpack.FileExists(dest); err != nil {
			return err
		} else if exists {
			if err := stager.LinkDirectoryInDepDir(dest, dirs[1]); err != nil {
				return err
			}
		}
	}
	for _, dirs := range [][]string{
		{"usr/lib/i386-linux-gnu/pkgconfig", "pkgconfig"},
		{"usr/lib/x86_64-linux-gnu/pkgconfig", "pkgconfig"},
		{"usr/lib/pkgconfig", "pkgconfig"},
		} {
			dest := filepath.Join(installDir, dirs[0])
			if exists, err := libbuildpack.FileExists(dest); err != nil {
				return err
			}else if exists {
				files, err := ioutil.ReadDir(dest)
				if err != nil {
					return err
				}
			destDir := filepath.Join(s.Stager.DepDir(), dirs[1])
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return err
			}
			for _, file := range files {
				//TODO: better way to copy a file?
				contents, err := ioutil.ReadFile(filepath.Join(dest, file.Name()))
				if err != nil {
					return err
				}
				newContents := strings.Replace(string(contents[:]), "prefix=/usr\n", "prefix="+filepath.Join(installDir, "usr")+"\n", -1)
				err = ioutil.WriteFile(filepath.Join(destDir, file.Name()), []byte(newContents), 0666)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil;
}
