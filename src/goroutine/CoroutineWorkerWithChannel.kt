import kotlinx.coroutines.*
import kotlinx.coroutines.channels.Channel
import kotlin.random.Random

fun main() = runBlocking {
    // Create a channel to communicate the results
    val results = Channel<Int>()
    val numWorkers = 5

    // Launch worker coroutines
    val jobs = List(numWorkers) { id ->
        launch {
            // Simulate some work by sleeping
            val sleepDuration = Random.nextInt(1, 6)
            println("Worker ${id + 1} is working for $sleepDuration seconds")
            delay(sleepDuration * 1000L)
            // Send the result to the channel
            results.send(sleepDuration)
        }
    }

    // Launch a coroutine to close the channel after all workers are done
    launch {
        jobs.forEach { it.join() }
        results.close()
    }

    // Collect all results from the channel
    for (result in results) {
        println("Received result: $result seconds")
    }

    println("All workers finished.")
}
