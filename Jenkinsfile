pipeline {
   agent any

    parameters {
        string defaultValue: 'api-server.app', description: 'Chatbot Server app name', name: 'server_app', trim: false
        string defaultValue: 'chatbot-cli.app', description: 'Chatbot CLI app name', name: 'cli_app', trim: false
        string defaultValue: 'iris-apis.image', description: 'Docker image name', name: 'docker_image', trim: false
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
            sh label: 'go unit test', script: 'go test'
        }
    }
    stage('build and archive executable') {
        steps {
            sh label: 'show version', script: 'go version'
            sh label: 'fetch swagger cmd', script: 'go get -u github.com/swaggo/swag/cmd/swag'
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
            sh label: 'Export docker image', script: "sudo docker save rayer/iris-apis > ${params.docker_image}"
            archiveArtifacts artifacts: "${params.docker_image}", fingerprint: true, followSymlinks: true, onlyIfSuccessful: true
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
            sh label: 'Kill image if exist', script: 'ssh jenkins@node.rayer.idv.tw docker kill APIService'
            sh label: 'Redeploy image', script: 'ssh jenkins@node.rayer.idv.tw docker run --rm --name APIService -p 8800:8080 -v ~/iris-apis:/app/config --hostname $(hostname) --rm -d rayer/iris-apis'
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
