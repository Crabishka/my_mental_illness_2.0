package ru.pryadchenko.stablepush.api

import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST
import retrofit2.http.Path

data class DeviceRequest(
    val model: String,
    val token: String
)

data class Device(
    val id: Int,
    val model: String,
    val token: String,
    val first_seen_at: String,
    val last_seen_at: String
)

data class NotificationRequest(
    val title: String,
    val body: String
)

interface DeviceApi {
    @POST("devices")
    suspend fun registerDevice(@Body request: DeviceRequest): Response<Device>
    
    @POST("notify/device/{id}")
    suspend fun sendNotification(
        @Path("id") deviceId: Int,
        @Body notification: NotificationRequest
    ): Response<Unit>
} 