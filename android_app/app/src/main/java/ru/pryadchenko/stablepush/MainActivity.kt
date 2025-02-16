package ru.pryadchenko.stablepush

import android.os.Bundle
import android.os.Build
import android.widget.Button
import android.widget.Toast
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.ViewCompat
import androidx.core.view.WindowInsetsCompat
import androidx.constraintlayout.widget.ConstraintLayout
import androidx.constraintlayout.widget.ConstraintSet
import com.google.firebase.messaging.FirebaseMessaging
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import ru.pryadchenko.stablepush.api.DeviceApi
import ru.pryadchenko.stablepush.api.DeviceRequest
import ru.pryadchenko.stablepush.api.NotificationRequest

class MainActivity : AppCompatActivity() {
    private val api = Retrofit.Builder()
        .baseUrl("https://4272517-lw36995.twc1.net/") // localhost для эмулятора
        .addConverterFactory(GsonConverterFactory.create())
        .build()
        .create(DeviceApi::class.java)

    private var deviceId: Int? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContentView(R.layout.activity_main)
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main)) { v, insets ->
            val systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars())
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom)
            insets
        }

        // Добавляем кнопку программно в центр экрана
        val button = Button(this).apply {
            text = "Отправить тестовый push"
            setOnClickListener { sendTestPush() }
        }
        
        // Создаем ConstraintLayout параметры для центрирования кнопки
        val params = ConstraintLayout.LayoutParams(
            ConstraintLayout.LayoutParams.WRAP_CONTENT,
            ConstraintLayout.LayoutParams.WRAP_CONTENT
        ).apply {
            topToTop = ConstraintLayout.LayoutParams.PARENT_ID
            bottomToBottom = ConstraintLayout.LayoutParams.PARENT_ID
            startToStart = ConstraintLayout.LayoutParams.PARENT_ID
            endToEnd = ConstraintLayout.LayoutParams.PARENT_ID
        }
        
        findViewById<ConstraintLayout>(R.id.main).addView(button, params)

        FirebaseMessaging.getInstance().token.addOnCompleteListener { task ->
            if (task.isSuccessful) {
                val token = task.result
                registerDevice(token)
            }
        }
    }

    private fun registerDevice(token: String) {
        CoroutineScope(Dispatchers.IO).launch {
            try {
                val response = api.registerDevice(
                    DeviceRequest(
                        model = "${Build.MANUFACTURER} ${Build.MODEL}",
                        token = token
                    )
                )
                if (response.isSuccessful) {
                    deviceId = response.body()?.id
                    println("Device registered successfully: ${response.body()}")
                } else {
                    println("Failed to register device: ${response.errorBody()?.string()}")
                }
            } catch (e: Exception) {
                println("Error registering device: ${e.message}")
            }
        }
    }

    private fun sendTestPush() {
        deviceId?.let { id ->
            CoroutineScope(Dispatchers.IO).launch {
                try {
                    val response = api.sendNotification(
                        deviceId = id,
                        notification = NotificationRequest(
                            title = "Тестовое уведомление",
                            body = "Привет от меня!"
                        )
                    )
                    if (response.isSuccessful) {
                        println("Notification sent successfully")
                    } else {
                        println("Failed to send notification: ${response.errorBody()?.string()}")
                    }
                } catch (e: Exception) {
                    println("Error sending notification: ${e.message}")
                }
            }
        } ?: run {
            Toast.makeText(this, "Устройство еще не зарегистрировано", Toast.LENGTH_SHORT).show()
        }
    }
}