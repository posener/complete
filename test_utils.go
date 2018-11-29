package complete

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TempMkdir creates a temporary directory
func TempMkdir(parentDir string, newDirPrefix string) (string, error) {
	parentDir = filepath.FromSlash(parentDir)
	dir, err := ioutil.TempDir(parentDir, newDirPrefix)
	if err != nil {
		return "", fmt.Errorf("failed to create dir with prefix %s in directory %s. Error %v", newDirPrefix, parentDir, err)
	}
	return dir, nil
}

// TempMkFile creates a temporary file.
func TempMkFile(dir string, fileName string) (string, error) {
	dir = filepath.FromSlash(dir)
	f, err := ioutil.TempFile(dir, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create test file %s in dir %s. Error %v", fileName, dir, err)
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// FileType custom type to indicate type of file
type FileType int

const (
	// RegularFile enum to represent regular file
	RegularFile FileType = 0
	// Directory enum to represent directory
	Directory FileType = 1
)

// ModificationType custom type to indicate file modification type
type ModificationType string

const (
	// UPDATE enum representing update operation on a file
	UPDATE ModificationType = "update"
	// CREATE enum representing create operation for a file/folder
	CREATE ModificationType = "create"
	// DELETE enum representing delete operation for a file/folder
	DELETE ModificationType = "delete"
	// APPEND enum representing append operation on a file
	APPEND ModificationType = "append"
)

// FileProperties to contain meta-data of a file like, file/folder name, file/folder parent dir, file type and desired file modification type
type FileProperties struct {
	FilePath         string
	FileParent       string
	FileType         FileType
	ModificationType ModificationType
}

// SimulateFileModifications mock function to simulate requested file/folder operation
// Parameters:
//	basePath: The parent directory for file/folder involved in desired file operation
//	fileModification: Meta-data of file/folder
// Returns:
//	path to file/folder involved in the operation
//	error if any or nil
func SimulateFileModifications(basePath string, fileModification FileProperties) (string, error) {
	// Files/folders intended to be directly under basepath will be indicated by fileModification.FileParent set to empty string
	if fileModification.FileParent != "" {
		// If fileModification.FileParent is not empty, use it to generate file/folder absolute path
		basePath = filepath.Join(basePath, fileModification.FileParent)
	}

	switch fileModification.ModificationType {
	case CREATE:
		if fileModification.FileType == Directory {
			filePath, err := TempMkdir(basePath, fileModification.FilePath)
			// t.Logf("In simulateFileModifications, Attempting to create folder %s in %s. Error : %v", fileModification.filePath, basePath, err)
			return filePath, err
		} else if fileModification.FileType == RegularFile {
			folderPath, err := TempMkFile(basePath, fileModification.FilePath)
			// t.Logf("In simulateFileModifications, Attempting to create file %s in %s", fileModification.filePath, basePath)
			return folderPath, err
		}
	case DELETE:
		if fileModification.FileType == Directory {
			return filepath.Join(basePath, fileModification.FilePath), os.RemoveAll(filepath.Join(basePath, fileModification.FilePath))
		} else if fileModification.FileType == RegularFile {
			return filepath.Join(basePath, fileModification.FilePath), os.Remove(filepath.Join(basePath, fileModification.FilePath))
		}
	case UPDATE:
		if fileModification.FileType == Directory {
			return "", fmt.Errorf("Updating directory %s is not supported", fileModification.FilePath)
		} else if fileModification.FileType == RegularFile {
			f, err := os.Open(filepath.Join(basePath, fileModification.FilePath))
			if err != nil {
				return "", err
			}
			if _, err := f.WriteString("Hello from Odo"); err != nil {
				return "", err
			}
			if err := f.Sync(); err != nil {
				return "", err
			}
			if err := f.Close(); err != nil {
				return "", err
			}
			return filepath.Join(basePath, fileModification.FilePath), nil
		}
	case APPEND:
		if fileModification.FileType == RegularFile {
			err := ioutil.WriteFile(filepath.Join(basePath, fileModification.FilePath), []byte("// Check watch command"), os.ModeAppend)
			if err != nil {
				return "", err
			}
			return filepath.Join(basePath, fileModification.FilePath), nil
		} else {
			return "", fmt.Errorf("Append not supported for file of type %v", fileModification.FileType)
		}
	default:
		return "", fmt.Errorf("Unsupported file operation %s", fileModification.ModificationType)
	}
	return "", nil
}

// CreateDirTree sets up a mock directory tree
// Parameters:
//  srcParentPath: The base path where src/dir tree is expected to be rooted
//	srcName: Name of the source directory
//	requiredFilePaths: list of required sources, their description like whether regularfile/directory, parent directory path of source and desired modification type like update/create/delete/append
// Returns:
//	absolute base path of source code
//	directory structure containing mappings from desired relative paths to their respective absolute path.
//  error if any
func CreateDirTree(srcParentPath string, srcName string, requiredFilePaths []FileProperties) (string, map[string]string, error) {

	if srcParentPath == `~` {
		srcParentPath = fixPathForm(srcParentPath, true, srcParentPath)
	}
	// This is required because ioutil#TempFile and ioutil#TempFolder creates paths with random numeric suffixes.
	// So, to be able to refer to the file/folder at any later point in time the created paths returned by ioutil#TempFile or ioutil#TempFolder will need to be saved.
	dirTreeMappings := make(map[string]string)

	// Create temporary directory for mock component source code
	srcPath, err := TempMkdir(srcParentPath, srcName)
	if err != nil {
		return "", dirTreeMappings, fmt.Errorf("failed to create dir %s under %s. Error: %v", srcName, srcParentPath, err)
	}
	dirTreeMappings[srcName] = srcPath

	// For each of the passed(desired) files/folders under component source
	for _, fileProperties := range requiredFilePaths {

		// get relative path using file parent and file name passed
		relativePath := filepath.Join(fileProperties.FileParent, fileProperties.FilePath)

		// get its absolute path using the mappings preserved from previous creates
		if realParentPath, ok := dirTreeMappings[fileProperties.FileParent]; ok {
			// real path for the intended file operation is obtained from previously maintained directory tree mappings by joining parent path and file name
			realPath := filepath.Join(realParentPath, fileProperties.FilePath)
			// Preserve the new paths for further reference
			fileProperties.FilePath = filepath.Base(realPath)
			fileProperties.FileParent, _ = filepath.Rel(srcPath, filepath.Dir(realPath))
		}

		// Perform mock operation as requested by the parameter
		newPath, err := SimulateFileModifications(srcPath, fileProperties)
		dirTreeMappings[relativePath] = newPath
		if err != nil {
			return "", dirTreeMappings, fmt.Errorf("unable to setup test env. Error %v", err)
		}

		fileProperties.FilePath = filepath.Base(newPath)
		fileProperties.FileParent = filepath.Dir(newPath)
	}

	// Return base source path and directory tree mappings
	return srcPath, dirTreeMappings, nil
}
