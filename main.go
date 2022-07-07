package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

const ItemFlag = "item"
const FileNameFlag = "fileName"
const OperationFlag = "operation"
const ListFlag = "list"
const IdFlag = "id"
const ADD = "add"
const LIST = "list"
const FindById = "findById"
const REMOVE = "remove"

func ValidateFileFlag(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}
	return nil
}

func ValidateItemFlag(data string) error {
	if data == "" {
		return fmt.Errorf("-item flag has to be specified")
	}
	return nil
}

func ValidateOperationFlag(operation string) error {

	if operation == "" {
		return fmt.Errorf("-operation flag has to be specified")
	} else if !(operation == ADD || operation == REMOVE || operation == LIST || operation == FindById) {
		return fmt.Errorf("Operation %s not allowed!", operation)
	}

	return nil

}
func Perform(args Arguments, writer io.Writer) error {

	err := ValidateOperationFlag(args[OperationFlag])

	if err != nil {
		return err
	}

	if ADD == args[OperationFlag] {

		if err := validateWriteMandatoryFlags(args); err != nil {
			return err
		}

		data, err1 := CombineDate(args[FileNameFlag], decode(args[ItemFlag]))

		if err1 != nil && !(strings.Contains(err1.Error(), "Item with id")) {
			writer.Write([]byte(err1.Error()))
			return err1
		} else if err1 != nil && strings.Contains(err1.Error(), "Item with id") {
			writer.Write([]byte(err1.Error()))
			return nil

		}

		WriteToFile(args[FileNameFlag], data)
	} else if FindById == args[OperationFlag] {
		data, err1 := findUserById(args)

		if err1 != nil && !(err1.Error() == "") {
			writer.Write([]byte(err1.Error()))
			return err1
		} else if err1 != nil && err1.Error() == "" {
			writer.Write([]byte(err1.Error()))
			return nil
		}

		out, _ := json.Marshal(data)
		writer.Write(out)

	} else if REMOVE == args[OperationFlag] {

		err1 := removeUserById(args)

		if err1 != nil && !(strings.Contains(err1.Error(), "Item with id")) {
			return err1
		} else if err1 != nil && strings.Contains(err1.Error(), "Item with id") {
			writer.Write([]byte(err1.Error()))
		}
	} else if LIST == args[OperationFlag] {

		data, err1 := listAllItems(args)

		if err1 != nil {
			return err1
		}

		if len(data) > 0 {

			out, _ := json.Marshal(data)
			writer.Write(out)
		}
	}
	return nil
}

func listAllItems(args Arguments) ([]Item, error) {

	err := ValidateFileFlag(args[FileNameFlag])

	if err != nil {
		return []Item{}, err
	}

	return ReadFromFile(args[FileNameFlag]), nil

}

func removeUserById(args Arguments) error {

	err := validateFindByIdFlags(args)

	if err != nil {
		return err
	}

	err1 := removeData(args, ReadFromFile(args[FileNameFlag]))

	if err1 != nil {
		return err1
	}

	return nil

}

func removeData(arguments Arguments, existingData []Item) error {

	out := make([]Item, 0)
	exist := false
	for _, v := range existingData {

		if v.Id != arguments[IdFlag] {
			out = append(out, v)
		} else {
			exist = true
		}
	}

	if exist {
		WriteToFile(arguments[FileNameFlag], out)
	} else {
		return fmt.Errorf("Item with id %s not found", arguments[IdFlag])
	}

	return nil

}

func findUserById(args Arguments) (Item, error) {

	err := validateFindByIdFlags(args)
	if err != nil {
		return Item{}, err
	}

	data := ReadFromFile(args[FileNameFlag])

	return filterSpecificUser(data, args[IdFlag])
}

func filterSpecificUser(data []Item, id string) (Item, error) {

	for _, v := range data {

		if v.Id == id {
			return v, nil
		}

	}

	return Item{}, fmt.Errorf("")

}

func validateFindByIdFlags(args Arguments) error {
	if err := ValidateFileFlag(args[FileNameFlag]); err != nil {
		return err
	} else if args[IdFlag] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	return nil

}

func validateWriteMandatoryFlags(args Arguments) error {

	if err := ValidateFileFlag(args[FileNameFlag]); err != nil {
		return err
	} else if err := ValidateItemFlag(args[ItemFlag]); err != nil {
		return err
	}

	return nil
}

func decode(data string) []Item {
	out := make([]Item, 0)
	item := Item{}
	errr := json.Unmarshal([]byte(data), &item)

	if errr != nil {
		fmt.Println(errr)
	}
	out = append(out, item)
	return out
}

func main() {
	err := Perform(parseArgs(), os.Stdout)

	if err != nil {
		panic(err)
	}
}
func CombineDate(fileName string, data []Item) ([]Item, error) {

	existingData := ReadFromFile(fileName)

	if hasDuplicateId(existingData, data[0]) {
		return nil, fmt.Errorf("Item with id %s already exists", data[0].Id)
	}
	return append(data, existingData...), nil
}

func hasDuplicateId(existingData []Item, dataToBeAdded Item) bool {

	for _, v := range existingData {

		if v.Id == dataToBeAdded.Id {
			return true
		}
	}

	return false

}
func WriteToFile(fileName string, data []Item) {

	file, _ := os.Create(fileName)
	marshalledData, _ := json.Marshal(data)
	defer file.Close()
	file.Write(marshalledData)

}

func ReadFromFile(fileName string) []Item {
	file1, _ := os.ReadFile(fileName)
	out := make([]Item, 0)
	json.Unmarshal(file1, &out)
	return out
}

func parseArgs() Arguments {
	operationFlag := flag.String(OperationFlag, "ADD", "supported operations: ADD remove etc")
	fileNameFlag := flag.String(FileNameFlag, "", "file name")
	itemFlag := flag.String(ItemFlag, "", "")
	idFlag := flag.String(IdFlag, "", "")

	flag.Parse()
	args := make(Arguments)

	args[FileNameFlag] = *fileNameFlag
	args[OperationFlag] = *operationFlag
	args[ItemFlag] = *itemFlag
	args[IdFlag] = *idFlag
	return args
}
