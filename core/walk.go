package core

import (
    "errors"
    "github.com/giwiro/walkline/utils"
    "io/fs"
    "io/ioutil"
    "os"
    "path/filepath"
    "sort"
)

type TransverseFunction func(node *MigrationNode) error

func processMigrationFile(path string, name string, version *Version, channel chan *MigrationFile) {
    file, err := ioutil.ReadFile(path)

    if err != nil {
        return
    }

    content := string(file)

    channel <- &MigrationFile{
        FilePath: path,
        FileName: name,
        Version:  version,
        Content:  content,
    }
}

func BuildMigrationTreeFromPath(dir string) (*MigrationNode, *[]*MigrationFailedFile, error) {
    workingDir, err := utils.GetWorkingDir()

    if err != nil {
        return nil, nil, err
    }

    tries := []string{
        filepath.Join(workingDir, "/src/main/resources/db/migrations/"),
    }

    if dir != "" {
        tries = []string{dir}
    }

    for _, d := range tries {
        if _, err := os.Stat(d); os.IsNotExist(err) {
            continue
        }

        firstNode, failedNodes, err := BuildMigrationTree(d)

        if err != nil {
            continue
        }

        return firstNode, failedNodes, nil
    }

    return nil, nil, errors.New("empty migration tree")
}

func BuildMigrationTree(dir string) (*MigrationNode, *[]*MigrationFailedFile, error) {
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        return nil, nil, errors.New("directory does not exist")
    }

    var failedFileNames []*MigrationFailedFile
    var files []*MigrationFile
    var filesLength = 0
    channel := make(chan *MigrationFile)

    err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if !info.IsDir() {
            version, err := ParseVersion(info.Name())
            if err != nil {
                failedFileNames = append(failedFileNames, &MigrationFailedFile{
                    FilePath: path,
                    FileName: info.Name(),
                    Error:    errors.New("failed file name validation"),
                })
                return errors.New("failed file name validation")
            }
            filesLength++
            go processMigrationFile(path, info.Name(), version, channel)
        }
        return nil
    })

    for i := 0; i < filesLength; i++ {
        files = append(files, <-channel)
    }

    if len(files) == 0 {
        return nil, nil, errors.New("no migration files")
    }

    sort.SliceStable(files, func(i, j int) bool {
        if files[i].Version.Version != files[j].Version.Version {
            return CompareVersionStr(files[i].Version.Version, files[j].Version.Version)
        }
        // V needs to be before U
        return files[i].Version.Prefix > files[j].Version.Prefix
    })

    if err != nil {
        return nil, nil, errors.New("could not read files from directory")
    }

    var firstNode = &MigrationNode{
        File:              files[0],
        UndoMigrationNode: nil,
        NextMigrationNode: nil,
        PrevMigrationNode: nil,
    }

    var iterNode = firstNode

    for i := 1; i < len(files); i++ {
        var node = &MigrationNode{
            File:              files[i],
            UndoMigrationNode: nil,
            NextMigrationNode: nil,
            PrevMigrationNode: iterNode,
        }

        if node.File.Version.Prefix == "U" {
            if iterNode.File.Version.Version != node.File.Version.Version {
                failedFileNames = append(failedFileNames, &MigrationFailedFile{
                    FilePath: files[i].FilePath,
                    FileName: files[i].FileName,
                    Error:    errors.New("undo migration version mismatch"),
                })
            }

            if iterNode.UndoMigrationNode != nil {
                failedFileNames = append(failedFileNames, &MigrationFailedFile{
                    FilePath: files[i].FilePath,
                    FileName: files[i].FileName,
                    Error:    errors.New("base version already got undo migration"),
                })
            }

            iterNode.UndoMigrationNode = node
        }

        if node.File.Version.Prefix == "V" {
            iterNode.NextMigrationNode = node
            iterNode = node
        }
    }

    return firstNode, &failedFileNames, nil
}

func TransverseMigrationTree(root *MigrationNode, fn TransverseFunction) {
    var iterNode = root

    for iterNode != nil {
        err := fn(iterNode)

        if err != nil {
            break
        }

        iterNode = iterNode.NextMigrationNode
    }
}

func FindMigrationNode(root *MigrationNode, version *VersionShort) *MigrationNode {
    var n *MigrationNode = nil
    TransverseMigrationTree(root, func(node *MigrationNode) error {
        if EqualsVersionFullAndShort(version, node.File.Version) {
            n = node
            return errors.New("found node")
        }
        return nil
    })
    return n
}
