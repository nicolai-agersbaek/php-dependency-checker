library("smartweb-libs")
properties([
        [$class: 'BuildDiscarderProperty', strategy: [$class: 'LogRotator', numToKeepStr: '50']]
    ]
)
node('worker'){
    withDockerRegistry([credentialsId: 'docker-hub']) {
        checkout scm
        sh 'docker build --force-rm --tag=smartweb/php-dependency-checker .'
        sh 'docker push smartweb/php-dependency-checker'
    }
}
