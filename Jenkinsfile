pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'make drunkenfall'
            }
        }
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Deploy') {
            steps {
                echo 'make deploy'
            }
        }
    }
}