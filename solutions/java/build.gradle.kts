plugins {
    java
}

java {
    toolchain {
        languageVersion = JavaLanguageVersion.of(21)
    }
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.openjdk.jmh:jmh-core:1.37")
    annotationProcessor("org.openjdk.jmh:jmh-generator-annprocess:1.37")
    implementation("org.postgresql:postgresql:42.7.4")
    implementation("com.uber:h3:4.1.1")
    implementation("com.google.geometry:s2-geometry:2.0.0")
    implementation("org.xerial:sqlite-jdbc:3.46.1.0")

    testImplementation(platform("org.junit:junit-bom:5.11.4"))
    testImplementation("org.junit.jupiter:junit-jupiter")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher")
}

tasks.withType<Test>().configureEach {
    useJUnitPlatform()
    systemProperty("labs.root", rootDir.resolve("../..").canonicalPath)
    testLogging {
        events("failed")
        exceptionFormat = org.gradle.api.tasks.testing.logging.TestExceptionFormat.SHORT
    }
}

tasks.test {
    exclude("**/*IntegrationTest*")
}

// Needs the disposable PostGIS from compose.yaml at the repository root.
val integrationTest by tasks.registering(Test::class) {
    description = "Runs the PostgreSQL/PostGIS integration tests."
    group = "verification"
    testClassesDirs = sourceSets.test.get().output.classesDirs
    classpath = sourceSets.test.get().runtimeClasspath
    include("**/*IntegrationTest*")
}
