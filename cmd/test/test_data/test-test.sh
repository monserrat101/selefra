#!/usr/bin/env bash
#######################################################################################################################
#                                                                                                                     #
#                              This script helps you test interactive programs                                        #
#                                                                                                                     #
#                                                                                                                     #
#                                                                                                   Version: 0.0.1    #
#                                                                                                                     #
#######################################################################################################################

# for command `selefra test`
cd ../../../
go build
rm ./cmd/test/test_data/module_test/selefra.exe
mv ./selefra.exe ./cmd/test/test_data/module_test/selefra.exe
cd ./cmd/test/test_data/module_test/
echo "begin run command selefra test"
./selefra.exe test $@

