package core

import (
	"errors"
	"fmt"
	"strconv"
)

type MigrationNode struct {
	File              *MigrationFile
	UndoMigrationNode *MigrationNode
	NextMigrationNode *MigrationNode
	PrevMigrationNode *MigrationNode
}

type MigrationFile struct {
	FilePath string
	FileName string
	Version  *Version
	Content  string
}

type MigrationFailedFile struct {
	FilePath string
	FileName string
	Error    error
}

func PrintMigrationTree(root *MigrationNode, currentVersion *VersionShort) {
	fmt.Println("\t[start]")
	TransverseMigrationTree(root, func(node *MigrationNode) error {
		var text string

		if currentVersion != nil && node.File.Version.Version == currentVersion.Version {
			text = "(curr)\t"
		} else {
			text = "\t"
		}

		text += fmt.Sprintf("%s (%s)", node.File.Version.Prefix + node.File.Version.Version, node.File.Version.Description)

		if node.UndoMigrationNode != nil {
			text += fmt.Sprintf("\t-> %s (%s)", node.UndoMigrationNode.File.Version.Prefix + node.UndoMigrationNode.File.Version.Version, node.UndoMigrationNode.File.Version.Description)
		}

		fmt.Println(text)
		return nil
	})
	fmt.Println("\t[end]")
}

func GenerateMigrationString(node *MigrationNode) string {
	var sql = "----------------- Migration " + node.File.Version.Prefix + node.File.Version.Version + " -----------------\n"
	sql += node.File.Content + "\n"
	sql += "--------------- End Migration " + node.File.Version.Prefix + node.File.Version.Version + " ---------------\n"
	sql += "\n"
	return sql
}

func GenerateMigrationStringFromVersionShortRange(flavor string, path string, schema string, currentVersion *VersionShort, leftVersion *VersionShort, rightVersion *VersionShort) (string, error) {
	var nodeList []*MigrationNode
	var migrationSqlString = ""
	var isSingleRevision = false

	if rightVersion != nil {
		isSingleRevision = CompareVersionShort(leftVersion, rightVersion)
	}

	firstNode, _, err := BuildMigrationTreeFromPath(path)
	var iterNode *MigrationNode

	if err != nil {
		return "", err
	}

	TransverseMigrationTree(firstNode, func(node *MigrationNode) error {
		if CompareVersionFullAndShort(leftVersion, node.File.Version) {
			nodeList = append(nodeList, node)
			iterNode = node.NextMigrationNode
			return errors.New("found first node")
		}

		if node.UndoMigrationNode != nil && CompareVersionFullAndShort(leftVersion, node.UndoMigrationNode.File.Version) {
			nodeList = append(nodeList, node.UndoMigrationNode)
			iterNode = node.UndoMigrationNode
			return errors.New("found first node")
		}
		return nil
	})

	if len(nodeList) == 0 {
		return "", errors.New("could not find first node")
	}

	if !isSingleRevision {
		TransverseMigrationTree(iterNode, func(node *MigrationNode) error {
			if node.File.Version.Prefix == "V" {
				nodeList = append(nodeList, node)

				if rightVersion != nil && CompareVersionFullAndShort(rightVersion, node.File.Version) {
					return errors.New("found second node")
				}
			}

			return nil
		})

		if rightVersion != nil {
			if len(nodeList) == 1 || !CompareVersionFullAndShort(rightVersion, nodeList[len(nodeList)-1].File.Version) {
				return "", errors.New("could not find last node")
			}
		}
	}

	for _, node := range nodeList {
		migrationSqlString += GenerateMigrationString(node)
	}

	migrationSqlString += GetInsertVersionQueryString(currentVersion, GetVersionShortFromFull(nodeList[len(nodeList)-1].File.Version), schema) + "\n"

	transaction, err := GenerateTransactionString(flavor, migrationSqlString)

	if err != nil {
		return "", err
	}

	return transaction, nil
}

func GenerateConsecutiveDowngradesMigrationString(flavor string, path string, schema string, currentVersion *VersionShort, times int) (string, error) {
	var nodeList []*MigrationNode
	var migrationSqlString = ""
	var iterNode *MigrationNode
	var iterTimes = times
	var finalVersion *VersionShort

	firstNode, _, err := BuildMigrationTreeFromPath(path)

	if err != nil {
		return "", err
	}

	var currentNode = FindMigrationNode(firstNode, currentVersion)

	if currentNode == nil {
		return "", err
	}

	iterNode = currentNode

	for iterNode != nil && iterTimes > 0 {
		if iterNode.UndoMigrationNode == nil {
			return "", errors.New("not enough consecutive undo migrations, " + strconv.Itoa(iterTimes) + " remaining")
		}

		nodeList = append(nodeList, iterNode.UndoMigrationNode)

		if iterNode.PrevMigrationNode != nil {
			finalVersion = GetVersionShortFromFull(iterNode.PrevMigrationNode.File.Version)
			iterNode = iterNode.PrevMigrationNode
		} else {
			finalVersion = nil
			iterNode = nil
		}
		iterTimes -= 1
	}

	if iterTimes > 0 {
		return "", errors.New("not enough consecutive undo migrations, " + strconv.Itoa(iterTimes) + " remaining")
	}

	if len(nodeList) == 0 {
		return "", errors.New("empty downgrades")
	}

	for _, node := range nodeList {
		migrationSqlString += GenerateMigrationString(node)
	}

	migrationSqlString += GetInsertVersionQueryString(currentVersion, finalVersion, schema)

	transaction, err := GenerateTransactionString(flavor, migrationSqlString)

	if err != nil {
		return "", err
	}

	return transaction, nil
}