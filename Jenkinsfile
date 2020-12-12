pipeline {
    agent any

    parameters {
        string defaultValue: 'api-server.app', description: 'Chatbot Server app name', name: 'server_app', trim: false
        string defaultValue: 'chatbot-cli.app', description: 'Chatbot CLI app name', name: 'cli_app', trim: false
    }

    stages {
        stage('Unit test') {
            steps {
                slackSend message: "${BUILD_TAG} start to build."
                sh label: 'go version', script: 'go version'
                sh label: 'install gocover-cobertura', script: 'go get github.com/t-yuki/gocover-cobertura'
                sh label: 'generate mocks', script: 'go generate ./...'
                sh label: 'go unit test', script: "FIXERIO_KEY=\"${FIXERIO_KEY}\" TEST_DB_CONN_STR=\"${TEST_DB_CONN_STR}\"; go test ./... --coverprofile=cover.out"
                sh label: 'convert coverage xml', script: '~/go/bin/gocover-cobertura < cover.out > coverage.xml'
            }
        }
        stage ('Extract test results') {
            steps {
                cobertura coberturaReportFile: 'coverage.xml'
            }
        }

        stage('Build and archive executable') {
            steps {
                sh label: 'show version', script: 'go version'
                sh label: 'generate documents', script: 'cd server && ~/go/bin/swag init -g server_main.go'
                sh label: 'build server', script: "cd server && go build -o ../bin/${params.server_app}"
                sh label: 'build cli', script: "cd cli && go build -o ../bin/${params.cli_app}"
                archiveArtifacts artifacts: 'bin/*', fingerprint: true, followSymlinks: false, onlyIfSuccessful: true
            }
        }
        stage('Build docker image') {
            steps {
                echo 'Building docker image'
                sh label: 'Build docker images', script: 'sudo docker build . -t rayer/iris-apis'
            }
        }
        stage('Push to docker repository') {
            steps {
                echo 'Pushing docker image'
                sh label: 'Push docker image', script: 'sudo docker push rayer/iris-apis'
            }
        }
        stage('Deploy image to api-test') {
            steps {
                echo 'Deploying docker image'
                sh label: 'Pull new images', script: 'ssh jenkins@node.rayer.idv.tw docker pull rayer/iris-apis:latest'
                catchError(buildResult: 'SUCCESS', stageResult: 'SUCCESS') {
                    sh label: 'Kill container if exist', script: 'ssh jenkins@node.rayer.idv.tw docker kill APIService-Test'
                }
                sh label: 'Redeploy container', script: 'ssh jenkins@node.rayer.idv.tw docker run --name APIService-Test -p 8801:8080 -v ~/iris-apis/test:/app/config -v /var/run/docker.sock:/var/run/docker.sock --hostname $(hostname) --rm -d rayer/iris-apis:latest'
            }
        }
        stage('Verify changes in test server') {
            steps {
                // The input step will prompt a message box for manual approval.
                slackSend message : "${BUILD_TAG} have been deployed to staging, please check ${BUILD_URL} for details."
                input message: "Deploy this version as release?"
                echo 'Verified: ${pwd()}'
            }
        }
        stage('Tag docker image as release') {
            steps {
                echo 'Pushing docker image'
                sh label: 'tag as release', script: 'sudo docker tag rayer/iris-apis:latest rayer/iris-apis:release'
                sh label: 'Push docker image', script: 'sudo docker push rayer/iris-apis:release'
            }
        }
        stage('Deploy image to api') {
            steps {
                echo 'Deploying docker image'
                sh label: 'Pull new images', script: 'ssh jenkins@node.rayer.idv.tw docker pull rayer/iris-apis:release'
                catchError(buildResult: 'SUCCESS', stageResult: 'SUCCESS') {
                    sh label: 'Kill container if exist', script: 'ssh jenkins@node.rayer.idv.tw docker kill APIService'
                }
                sh label: 'Redeploy container', script: 'ssh jenkins@node.rayer.idv.tw docker run --name APIService -p 8800:8080 -v ~/iris-apis:/app/config -v /var/run/docker.sock:/var/run/docker.sock --hostname $(hostname) --rm -d rayer/iris-apis:release'
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
