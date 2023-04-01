#!/bin/bash

FEED_URL="https://asobistore.jp/special/List?tag_seq%5B%5D=1"
ASOBI_STORE="https://asobistore.jp/"


  # getSearchTerm "asobiURL"
function getSearchTermAsobi {
  local asobiUrl="$1"
  if [[ "$asobiUrl" =~ (List) ]]; then
    local episodeUrl="${asobi_store}$(curl -s "$asobiUrl" | xmllint --html --xpath "string(//ul[@class='list-main-product']/li/*/div[p='視聴制限なし']/../@href)" - 2>/dev/null)"
  else
    local episodeUrl="$asobiUrl"
  fi
  echo "$(curl -s "$episodeUrl"| xmllint --html --xpath "//ul[@class='list-dcm']/li[2]/text()" - 2>/dev/null)"
}

MY_ID="$(getSearchTermAsobi "$FEED_URL")"

echo $MY_ID