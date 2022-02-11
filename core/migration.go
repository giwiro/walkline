package core

import (
	"errors"
	"fmt"
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

func PrintMigrationTree(root *MigrationNode, currentVersion string) {
	fmt.Println("\t[start]")
	TransverseMigrationTree(root, func(node *MigrationNode) error {
		var text string

		if node.File.Version.Version == currentVersion {
			text = "(curr)\t"
		} else {
			text = "\t"
		}

		text += node.File.Version.Prefix + node.File.Version.Version + " (" + node.File.Version.Description + ")"

		if node.UndoMigrationNode != nil {
			text += "\t-> " + node.UndoMigrationNode.File.Version.Prefix + node.UndoMigrationNode.File.Version.Version + " (" + node.UndoMigrationNode.File.Version.Description + ")"
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

func GenerateMigrationStringFromVersionShortRange(flavor string, leftVersion *VersionShort, rightVersion *VersionShort) (string, error) {
	var nodeList []*MigrationNode
	var migrationSqlString = ""
	var isSingleRevision = CompareVersionShort(leftVersion, rightVersion)

	firstNode, _, err := BuildMigrationTree("/tmp/migrations")
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

				if CompareVersionFullAndShort(rightVersion, node.File.Version) {
					return errors.New("found second node")
				}
			}

			return nil
		})

		if len(nodeList) == 1 || !CompareVersionFullAndShort(rightVersion, nodeList[len(nodeList) - 1].File.Version) {
			return "", errors.New("could not find last node")
		}
	}

	for _, node := range nodeList {
		migrationSqlString += GenerateMigrationString(node)
	}

	migrationSqlString += GetSetDatabaseVersionQueryString(nodeList[len(nodeList)-1].File.Version) + "\n"

	transaction, err := GenerateTransactionString(flavor, migrationSqlString)

	if err != nil {
		return "", err
	}

	return transaction, nil
}
