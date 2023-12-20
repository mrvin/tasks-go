#!/bin/bash

get_film_info="./bin/get-film-info"
file_path_for_saving="movie_information.json"

reg_exp_year='^[0-9]{4}$'
reg_exp_api_key='^[a-z0-9]{8}$'

if [[ $# -ne 4 ]]; then
  echo -e "No required number of parameters."
  exit 1
fi

while getopts "f:k:" opt
do
  case $opt in
    f) if ! test -f "$OPTARG"; then
        echo -e "$OPTARG file not exist."
        exit 1
      fi
      input_file=$OPTARG;;
    k) if ! [[ $OPTARG =~ $reg_exp_api_key ]] ; then
        echo -e "$OPTARG not match api key."
        exit 1
      fi
      api_key=$OPTARG;;
    *) echo -e "No reasonable options found!"
      exit 1;;
  esac
done

while read -r line
do
  line="${line%.*}"
  film="${line%_*}"
  year="${line##*_}"

  if [[ $year =~ $reg_exp_year ]] ; then
    $get_film_info -k "$api_key" -n "$film" -y "$year" -f "$file_path_for_saving" -p 1>/dev/null
  else
    $get_film_info -k "$api_key" -n "$line" -f "$file_path_for_saving" -p 1>/dev/null
  fi
done < "$input_file"

exit 0
