ROOT_DIR=$PWD

LEGO_PACKAGES_PATH="acme/api/internal,acme,certcrypto,certificate,cmd,lego"
LEGO_PACKAGES_NAME="secure,api,0,0,0,0"

PEBBLE_PACKAGES_PATH="ca,cmd,ocsp,wfe"
PEBBLE_PACKAGES_NAME="0,pebble,0,0"

PROJECTS_TO_DOCUMENT=(go-lego go-pebble)
MODULES_NAME=(github.com/go-acme/lego/v4 github.com/letsencrypt/pebble/v2)
PROJECTS_PACKAGES_PATH=(${LEGO_PACKAGES_PATH} ${PEBBLE_PACKAGES_PATH})
PROJECTS_PACKAGES_NAME=(${LEGO_PACKAGES_NAME} ${PEBBLE_PACKAGES_NAME})

GENERATED_DOCS=${PWD}/temp
PUBLISH_DOCS=${PWD}/documentation

mkdir -p ${GENERATED_DOCS}

rm -rf ${PUBLISH_DOCS}/{pkg,src}
mkdir ${PUBLISH_DOCS}/{pkg,src}

j=0
for DIR in ${PROJECTS_TO_DOCUMENT[*]}
do
  MODULE_NAME=${MODULES_NAME[${j}]}
  
  IFS=',' read -r -a PACKAGES_PATH <<< "${PROJECTS_PACKAGES_PATH[${j}]}"
  IFS=',' read -r -a PACKAGES_NAME <<< "${PROJECTS_PACKAGES_NAME[${j}]}"  

  cd ${GENERATED_DOCS}
  cp -R {css,jvs,png} ${PUBLISH_DOCS}

  i=0
  for PACKAGE_PATH in ${PACKAGES_PATH[*]}
  do
    
    FULL_PACKAGE_NAME=a
    if [ "${PACKAGES_NAME[${i}]}" == "0" ]; then
      FULL_PACKAGE_NAME=${PACKAGE_PATH}
    else
      PACKAGE_NAME=${PACKAGES_NAME[${i}]}
      FULL_PACKAGE_NAME=${PACKAGE_PATH}/${PACKAGE_NAME}
    fi

    echo ${FULL_PACKAGE_NAME}  

    cd ${ROOT_DIR}/${DIR}/${FULL_PACKAGE_NAME}

    ~/go/bin/golds -gen -dir ${GENERATED_DOCS} -only-list-exporteds=false

    cd ${GENERATED_DOCS}

    if [ "${PACKAGES_NAME[${i}]}" == "0" ]; then
      mkdir -p ${PUBLISH_DOCS}/pkg/${MODULE_NAME}
    else
      mkdir -p ${PUBLISH_DOCS}/pkg/${MODULE_NAME}/${PACKAGE_PATH}
    fi
    
    cp pkg/${MODULE_NAME}/${FULL_PACKAGE_NAME}.html ${PUBLISH_DOCS}/pkg/${MODULE_NAME}/${FULL_PACKAGE_NAME}.html

    mkdir -p ${PUBLISH_DOCS}/src/${MODULE_NAME}/${FULL_PACKAGE_NAME}
    cp src/${MODULE_NAME}/${FULL_PACKAGE_NAME}/* ${PUBLISH_DOCS}/src/${MODULE_NAME}/${FULL_PACKAGE_NAME}

    ((i+=1))
  done
  ((j+=1))
done

# Generating documentation for go-tests

cd ${ROOT_DIR}/tests/src
~/go/bin/golds -gen -dir ${GENERATED_DOCS} -only-list-exporteds=false
cd ${GENERATED_DOCS}

mkdir -p ${PUBLISH_DOCS}/pkg/tls_tests
echo ${GENERATED_DOCS}
cp pkg/tls_tests/src.html ${PUBLISH_DOCS}/pkg/tls_tests/src.html

mkdir -p ${PUBLISH_DOCS}/src/tls_tests/src
cp src/tls_tests/src/* ${PUBLISH_DOCS}/src/tls_tests/src


# Generating documentation for go

STD_PACKAGES_NAME=(liboqs_sig tls x509)

cd ${ROOT_DIR}/go-lego/acme/api/internal/nonces

~/go/bin/golds -gen -dir ${GENERATED_DOCS} -only-list-exporteds=false

cd ${GENERATED_DOCS}

mkdir -p ${PUBLISH_DOCS}/pkg/crypto

for PACKAGE_NAME in ${STD_PACKAGES_NAME[*]}
do
  cp pkg/crypto/${PACKAGE_NAME}.html ${PUBLISH_DOCS}/pkg/crypto/${PACKAGE_NAME}.html

  mkdir -p ${PUBLISH_DOCS}/src/crypto/${PACKAGE_NAME}
  cp src/crypto/${PACKAGE_NAME}/* ${PUBLISH_DOCS}/src/crypto/${PACKAGE_NAME}
done

rm -rf ${GENERATED_DOCS}