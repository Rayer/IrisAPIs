pipeline {
    agent any

    stages {
        stage('Unit test') {
            steps {

                slackSend message: "${BUILD_TAG} start to build."
                sh label: 'go version', script: 'go version'
                sh label: 'install gocover-cobertura', script: 'go get github.com/t-yuki/gocover-cobertura'
                sh label: 'generate mocks', script: 'go generate ./...'
                lock(quantity: 1, resource: 'IrisAPI_UT') {
                    withCredentials([string(credentialsId: 'fixerioApiKey', variable: 'FIXERIO_KEY'), string(credentialsId: 'testConnectionString', variable: 'TEST_DB_CONN_STR')]) {
                        sh label: 'go unit test', script: "FIXERIO_KEY=\"${FIXERIO_KEY}\" TEST_DB_CONN_STR=\"${TEST_DB_CONN_STR}\"; go test ./... --coverprofile=cover.out"
                    }
                }
                sh label: 'convert coverage xml', script: '~/go/bin/gocover-cobertura < cover.out > coverage.xml'
            }
        }
        stage ('Extract test results') {
            steps {
                cobertura coberturaReportFile: 'coverage.xml'
            }
        }

        stage('Build docker image') {
            steps {
                echo 'Building docker image'
                sh label: 'Build docker images', script: "sudo docker build . --build-arg IMAGE_TAG=${BRANCH_NAME}-${BUILD_NUMBER} --build-arg JENKINS_LINK=${BUILD_URL} -t rayer/iris-apis:${BRANCH_NAME}-${BUILD_NUMBER}"
            }
        }
        stage('Push to docker repository') {
            steps {
                echo 'Pushing docker image'
                sh label: 'Push docker image', script: "sudo docker push rayer/iris-apis:${BRANCH_NAME}-${BUILD_NUMBER}"
            }
        }
        stage('Deploy image to api-test') {
            when {
                not {
                    branch "master"
                }
            }

            steps {
                lock(quantity: 1, resource: 'IrisAPITest_Deploying') {
                    echo 'Deploying docker image'
                    sh label: 'Pull new images', script: 'ssh jenkins@node.rayer.idv.tw docker pull rayer/iris-apis:${BRANCH_NAME}-${BUILD_NUMBER}'
                    catchError(buildResult: 'SUCCESS', stageResult: 'SUCCESS') {
                        sh label: 'Kill container if exist', script: 'ssh jenkins@node.rayer.idv.tw docker kill APIService-Test'
                    }
                    sh label: 'Redeploy container', script: 'ssh jenkins@node.rayer.idv.tw docker run --name APIService-Test -p 8801:8080 -p 9002:8082 -v ~/iris-apis/test:/app/config -v /var/run/docker.sock:/var/run/docker.sock --hostname $(hostname) --rm -d rayer/iris-apis:${BRANCH_NAME}-${BUILD_NUMBER}'
                }
            }
        }

        stage('Deploy image to api') {
            when {
                branch "master"
            }
            steps {
                echo 'Deploying docker image'
                sh label: 'Pull new images', script: 'ssh jenkins@node.rayer.idv.tw docker pull rayer/iris-apis:${BRANCH_NAME}-${BUILD_NUMBER}'
                catchError(buildResult: 'SUCCESS', stageResult: 'SUCCESS') {
                    sh label: 'Kill container if exist', script: 'ssh jenkins@node.rayer.idv.tw docker kill APIService'
                }
                sh label: 'Redeploy container', script: 'ssh jenkins@node.rayer.idv.tw docker run --name APIService -p 8800:8080 -p 9001:8082 -v ~/iris-apis:/app/config -v /var/run/docker.sock:/var/run/docker.sock --hostname $(hostname) --rm -d rayer/iris-apis:${BRANCH_NAME}-${BUILD_NUMBER}'
            }
        }
    }

    post {
        aborted {
            slackSend message: "${BUILD_TAG} build ended."
        }
        success {
            slackSend message: "${BUILD_TAG} is built successfully."
        }
        failure {
            slackSend message: "${BUILD_TAG} is failed to be built."
        }
    }
}
