@Library('jenkins-shared-lib') _

def pipelineConfig = [
    appName: 'Clementineqq/voidsounds',
    registry: 'ghcr.io'
]

def ciPipeline = load('jenkins-pipelines/ci/voidsounds.groovy')
def cdStaging = load('jenkins-pipelines/cd/staging.groovy')

node {
    checkout scm

    stage('CI') {
        ciPipeline.call(pipelineConfig)
    }

    stage('CD Staging') {
        cdStaging.call(pipelineConfig + [imageTag: env.BUILD_NUMBER])
    }
}