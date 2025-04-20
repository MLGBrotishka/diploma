package logger

var packageLogger *Logger

// Init - инициализация глобального логгера.
func Init(level string) {
	packageLogger = New(level)
}

// Debug - логгирование в debug уровне.
func Debug(message interface{}, args ...interface{}) {
	if packageLogger == nil {
		packageLogger = New(DebugLevel)
	}
	packageLogger.Debug(message, args...)
}

// Info - логгирование в info уровне.
func Info(message string, args ...interface{}) {
	if packageLogger == nil {
		packageLogger = New(DebugLevel)
	}
	packageLogger.Info(message, args...)
}

// Warn - логгирование в warn уровне.
func Warn(message string, args ...interface{}) {
	if packageLogger == nil {
		packageLogger = New(DebugLevel)
	}
	packageLogger.Warn(message, args...)
}

// Error - логгирование в error уровне.
func Error(message interface{}, args ...interface{}) {
	if packageLogger == nil {
		packageLogger = New(DebugLevel)
	}
	packageLogger.Error(message, args...)
}

// Fatal - логгирование в fatal уровне.
func Fatal(message interface{}, args ...interface{}) {
	if packageLogger == nil {
		packageLogger = New(DebugLevel)
	}
	packageLogger.Fatal(message, args...)
}
