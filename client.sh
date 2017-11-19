curl -X PUT -d '[{"key":{"encoding":"string","data":"key1"},"value":{"encoding":"string","data":"value1"}},{"key":{"encoding":"binary","data":"001011010"},"value":{"encoding":"string","data":"001010110"}}]' http://localhost:8080/set && printf "\n"
sleep 1

curl -X POST -d '[{"encoding":"binary","data":"001011010"}]' http://localhost:8080/fetch && printf "\n"
sleep 1

curl -X POST -d '[{"encoding":"string","data":"key1"},{"encoding":"string","data":"key2"}]' http://localhost:8080/fetch && printf "\n"
sleep 1

curl -X POST -d '[{"encoding":"string","data":"key1"},{"encoding":"binary","data":"001011010"}]' http://localhost:8080/fetch && printf "\n"
sleep 1

curl -X GET http://localhost:8080/fetch && printf "\n"
sleep 1

curl -X POST -d '[{"encoding":"string","data":"key1"},{"encoding":"string","data":"key2"}]' http://localhost:8080/query 


