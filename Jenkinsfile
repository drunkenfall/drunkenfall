pipeline {
  agent any

  stages {
    // Setup steps
    stage('Go dependencies') {
      steps {
        sh 'make download'
      }
    }
    stage('npm dependencies') {
      steps {
        sh 'make npm'
      }
    }

    // Verification steps
    stage('Verify') {
      steps {
        sh 'make test'
      }
      steps {
        sh 'make lint'
      }
      steps {
        sh 'make race'
      }
    }

    // Building step
    stage('Compile application') {
      steps {
        sh 'make drunkenfall'
      }
    }
    stage('Build npm') {
      steps {
        sh 'make npm-dist'
      }
    }

    // Deploy step
    stage('Deploy') {
      steps {
        echo 'make deploy'
      }
    }
  }
}