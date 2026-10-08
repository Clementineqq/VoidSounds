@Library('jenkins-shared-lib@dev') _

pipeline {
    agent any //встроенный узел дженинкс, там должны быть docker CLI и доступ к демону

    environment {
        APP_NAME  = 'clementineqq/voidsounds'
        REGISTRY  = 'ghcr.io'
        IMAGE_TAG = "${env.GIT_COMMIT.take(7)}" //короткий SHA коммита (первые 7 символов SHA-коммита)
    }

    options {
        disableConcurrentBuilds() //не запускать две сборки одной джобы одновременно

        
        buildDiscarder(logRotator(numToKeepStr: '10')) // хранить логи только 10 последних сборок
    }

    stages {
        stage('Test') { //проверка кода: go mod download + go vet + go test (внутри контейнера го)
            steps {
                runTests(action: 'all')
            }
        }

        stage('Build Image') { // сборка докер-образа из dockerfile приложения

            steps {
                buildImage(
                    imageName: "${APP_NAME}",
                    imageTag:  "${IMAGE_TAG}"
                )
            }
        }

    
        stage('Push to Registry') {
            steps {
                pushImage(
                    registry:      "${REGISTRY}",
                    imageName:     "${APP_NAME}",
                    imageTag:      "${IMAGE_TAG}",
                    credentialsId: 'github-registry'
                )
            }
        }
    }

    post {
        success {
            notify(status: 'success', message: "Built ${REGISTRY}/${APP_NAME}:${IMAGE_TAG}")
            build job: 'voidsounds-cd-staging', // запускаю CD-джобу после успешного пуша CI
                  parameters: [string(name: 'IMAGE_TAG', value: "${IMAGE_TAG}")], // прокидываю SHA, он попадает в params.IMAGE_TAG, оттуда в композ и в докер пулл
                  wait: false //жду деплой
        }

        failure {
            notify(status: 'failure', message: "CI failed for ${APP_NAME}")
        }
        always {
            echo "Сборка #${env.BUILD_NUMBER} завершена: ${currentBuild.currentResult}"
        }
    }
}