pipeline {
   agent any

    parameters {
        string defaultValue: 'api-server.app', description: 'Chatbot Server app name', name: 'server_app', trim: false
        string defaultValue: 'chatbot-cli.app', description: 'Chatbot CLI app name', name: 'cli_app', trim: false
    }

   stages {
    stage('Fetch from github') {
        steps {
            slackSend message: 'Project rayer/iris-apis start to build.'
            git credentialsId: '26c5c0a0-d02d-4d77-af28-761ffb97c5cc', url: 'https://github.com/Rayer/IrisAPIs.git'
        }
    }
    stage('Unit test') {
        steps {
            sh label: 'go version', script: 'go version'
            sh label: 'install gocover-cobertura', script: 'go get github.com/t-yuki/gocover-cobertura'
            sh label: 'go unit test', script: 'go test --coverprofile=cover.out'
            sh label: 'convert coverage xml', script: '~/go/bin/gocover-cobertura < cover.out > coverage.xml'
        }
    }
    stage ("Extract test results") {
        steps {
            cobertura coberturaReportFile: 'coverage.xml'
        }
    }

    stage('build and archive executable') {
        steps {
            sh label: 'show version', script: 'go version'
            sh label: 'generate documents', script: "cd server && ~/go/bin/swag init -g server_main.go"
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
            sh label: 'Push docker image', script: "sudo docker push rayer/iris-apis"
        }
      }
    stage('Deploy to target servers') {
        steps {
            echo 'Deploying docker image'
            sh label: 'Pull new images', script: 'ssh jenkins@node.rayer.idv.tw docker pull rayer/iris-apis'
            cacheError {
                sh label: 'Kill container if exist', script: 'ssh jenkins@node.rayer.idv.tw docker kill APIService'
            }
            sh label: 'Redeploy container', script: 'ssh jenkins@node.rayer.idv.tw docker run --name APIService -p 8800:8080 -v ~/iris-apis:/app/config --hostname $(hostname) --rm -d rayer/iris-apis'
        }
    }
   }

   post {
        aborted {
            slackSend message: 'Project rayer/iris-apis aborted.'
        }
        success {
            slackSend message: 'Project rayer/iris-apis is built successfully.'
        }
        failure {
            slackSend message: 'Project rayer/iris-apis is failed to be built.'
        }
    }
}
