pipeline {
    agent any

    environment {
        PROJECT_NAME = "platform"
        PROJECT_ENV = "dev"

        APPLICATION_NAME = "pki-server"
        APPLICATION_REPLICAS = 2

        REPOSITORY_AUTH = "github"
        REPOSITORY_URL = "https://github.com/buxiaomo/pki-server.git"

        REGISTRY_AUTH = "harbor"
        REGISTRY_HOST = "harbor.nuomitech.cn"
        REGISTRY_REPO = "platform"
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '15'))
    }

    stages {
        stage('checkout') {
            steps {
                checkout([$class: 'GitSCM', branches: [[name: '*/main']], extensions: [], userRemoteConfigs: [[credentialsId: "${env.REPOSITORY_AUTH}", url: "${env.REPOSITORY_URL}"]]])
            }
        }

        stage('build') {
            parallel {
                stage('Docker Image') {
                    steps{
                        sh label: '', script: "docker build -t ${env.REGISTRY_HOST}/${env.REGISTRY_REPO}/${env.APPLICATION_NAME}:${BUILD_ID} . -f Dockerfile"
                    }
                }
            }
        }

        stage('push') {
            steps {
                retry(3) {
                    script {
                        if (env.REGISTRY_AUTH) {
                           withDockerRegistry(credentialsId: "${env.REGISTRY_AUTH}", url: "https://${env.REGISTRY_HOST}") {
                                sh label: 'Docker', script: "docker push ${env.REGISTRY_HOST}/${env.REGISTRY_REPO}/${env.APPLICATION_NAME}:${BUILD_ID}"
                           }
                        } else {
                            sh label: 'Docker', script: "docker push ${env.REGISTRY_HOST}/${env.REGISTRY_REPO}/${env.APPLICATION_NAME}:${BUILD_ID}"
                        }
                    }
                }
            }
        }

        stage('deploy') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'harbor', passwordVariable: 'password', usernameVariable: 'username')]) {
                        sh label: 'add repo', script: "helm repo add --username $username --password $password platform https://${env.REGISTRY_HOST}/chartrepo/platform"
                        sh label: 'install chart', script: "helm upgrade -i ${env.APPLICATION_NAME} --username $username --password $password --version 0.1.0 platform/plaform -n ${env.PROJECT_NAME}-${env.PROJECT_ENV} --set readinessPath='/v1/pki/healthz' --set Image=${env.REGISTRY_HOST}/${env.REGISTRY_REPO}/${env.APPLICATION_NAME}:${BUILD_ID} --set Project=${env.PROJECT_NAME} --set Env=${env.PROJECT_ENV} --set Replicas=${env.APPLICATION_REPLICAS} --set Port=8080"
                    }
                }
            }
        }
    }
}
