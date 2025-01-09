# Go Concurrency Fundamentals

## 1. Goroutine
Goroutine คือ lightweight thread ใน Go ที่ทำให้เราสามารถรันโค้ดแบบ concurrent ได้

```go
// ตัวอย่างการใช้ goroutine
func main() {
    // รัน function ใน goroutine
    go func() {
        fmt.Println("Running in goroutine")
    }()

    // รัน function ที่มีอยู่แล้วใน goroutine
    go doSomething()

    // รอให้ goroutines ทำงานเสร็จ
    time.Sleep(1 * time.Second)
}
```

## 2. Channel
Channel คือท่อส่งข้อมูลระหว่าง goroutines

```go
func main() {
    // สร้าง buffered channel ขนาด 2
    ch := make(chan string, 2)

    go func() {
        ch <- "hello"  // ส่งข้อมูล
        ch <- "world"
        close(ch)      // ปิด channel เมื่อส่งเสร็จ
    }()

    // รับข้อมูลจาก channel
    for msg := range ch {
        fmt.Println(msg)
    }
}
```

## 3. Race Condition
Race Condition เกิดเมื่อ goroutines หลายตัวเข้าถึงข้อมูลเดียวกันพร้อมกัน

```go
func main() {
    counter := 0

    // เกิด race condition
    for i := 0; i < 1000; i++ {
        go func() {
            counter++ // อาจได้ผลลัพธ์ไม่ถูกต้อง
        }()
    }
}
```

## 4. Mutex
Mutex ใช้ล็อคการเข้าถึงข้อมูลที่ใช้ร่วมกัน

```go
func main() {
    counter := 0              // ตัวนับ
    var mutex sync.Mutex      // กุญแจล็อค
    var wg sync.WaitGroup     // ตัวนับงาน

    // สร้าง goroutines 1000 ตัว
    for i := 0; i < 1000; i++ {
        wg.Add(1)             // เพิ่มงาน
        go func() {
            defer wg.Done()    // เสร็จแล้วลดงาน
            
            mutex.Lock()       // ล็อคก่อนแก้ counter
            counter++          // เพิ่มค่า counter
            mutex.Unlock()     // ปลดล็อค
        }()
    }

    wg.Wait()                 // รอให้ทุก goroutine เสร็จ
    fmt.Println(counter)      // พิมพ์ผลลัพธ์
}
```

ถ้าไม่ใช้ mutex และ WaitGroup:
```go
func main() {
    counter := 0
    
    // BAD: ไม่มีการล็อค และไม่รู้ว่าทำงานเสร็จเมื่อไหร่
    for i := 0; i < 1000; i++ {
        go func() {
            counter++
        }()
    }
    
    // อาจจะพิมพ์ก่อน goroutines ทำงานเสร็จ
    // และ counter อาจไม่ถึง 1000
    fmt.Println(counter)
}
```

## 5. Atomic Operations
Atomic Operations ใช้สำหรับ operations ง่ายๆ ที่ต้องการความปลอดภัย

```go
func main() {
    var counter atomic.Int64
    var wg sync.WaitGroup

    // ใช้ atomic operations
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Add(1)
        }()
    }

    wg.Wait()
    fmt.Println("Counter:", counter.Load()) // จะได้ 1000 เสมอ
}
```

## 6. Select ใน Go

## 1. Select คืออะไร?
Select เป็นคำสั่งพิเศษใน Go ที่ใช้จัดการหลาย channels พร้อมกัน เหมือนเป็น "สวิตช์บอร์ด" ที่คอยดูว่า channel ไหนพร้อมทำงาน

## 2. รูปแบบการใช้ Select

### 2.1 แบบพื้นฐาน:
```go
select {
case msg := <-ch1:
    // ทำงานเมื่อได้ข้อมูลจาก ch1
case ch2 <- data:
    // ทำงานเมื่อส่งข้อมูลไป ch2 สำเร็จ
default:
    // ทำงานเมื่อไม่มี case ไหนพร้อม
}
```

### 2.2 แบบมี Timeout:
```go
select {
case msg := <-ch:
    fmt.Println("ได้ข้อมูล:", msg)
case <-time.After(1 * time.Second):
    fmt.Println("รอนานเกินไป!")
}
```

### 2.3 แบบไม่ Block:
```go
select {
case ch <- data:
    fmt.Println("ส่งสำเร็จ")
default:
    fmt.Println("channel ไม่พร้อม ไม่รอ")
}
```

## 3. การทำงานของ Select

### 3.1 กรณีไม่มี Default:
- รอจนกว่าจะมี case ที่พร้อมทำงาน
- ถ้ามีหลาย case พร้อมพร้อมกัน จะสุ่มเลือก
```go
select {
case <-ch1:
    fmt.Println("ได้จาก ch1")
case <-ch2:
    fmt.Println("ได้จาก ch2")
}
// จะ block จนกว่า ch1 หรือ ch2 จะพร้อม
```

### 3.2 กรณีมี Default:
- ถ้าไม่มี case ไหนพร้อม จะทำ default ทันที
- ไม่ block
```go
select {
case <-ch:
    fmt.Println("ได้ข้อมูล")
default:
    fmt.Println("ไม่มีข้อมูล ไปทำงานอื่นต่อ")
}
```

## 4. การใช้ Select กับ For

### 4.1 แบบทำครั้งเดียว:
```go
// เหมาะกับการเช็คครั้งเดียวแล้วไปต่อ
select {
case data := <-ch:
    process(data)
case <-time.After(1 * time.Second):
    handleTimeout()
}
```

### 4.2 แบบทำซ้ำ:
```go
// เหมาะกับ worker ที่ต้องทำงานต่อเนื่อง
for {
    select {
    case data := <-ch:
        process(data)
    case <-done:
        return
    }
}
```

## 5. ตัวอย่างการใช้งานจริง

### 5.1 Worker Pool:
```go
func worker(jobs <-chan int, results chan<- int, done <-chan bool) {
    for {
        select {
        case job := <-jobs:
            results <- process(job)
        case <-done:
            return
        }
    }
}
```

### 5.2 Rate Limiting:
```go
func rateLimiter(requests <-chan Request, limit time.Duration) {
    ticker := time.NewTicker(limit)
    for {
        select {
        case req := <-requests:
            process(req)
        case <-ticker.C:
            // รอจนกว่าจะถึงเวลา
        }
    }
}
```

### 5.3 Graceful Shutdown:
```go
func (s *Server) Run() {
    for {
        select {
        case req := <-s.requests:
            handleRequest(req)
        case <-s.shutdown:
            // 1. หยุดรับ requests ใหม่
            // 2. รอให้ requests ที่กำลังทำงานเสร็จ
            // 3. ปิดระบบ
            return
        }
    }
}
```

## 6. ข้อควรระวัง

1. Select แบบไม่มี case:
```go
select {}  // deadlock ทันที!
```

2. Select ที่มีแต่ default:
```go
// ไม่มีประโยชน์ เพราะจะทำ default ทันที
select {
default:
    fmt.Println("ทำงานทันที")
}
```

3. การสุ่มเลือก case:
```go
// ระวัง! ถ้ามีหลาย case พร้อมกัน
// จะสุ่มเลือก ไม่ทำตามลำดับ
select {
case <-ch1:  // อาจไม่ได้ทำก่อน
case <-ch2:  // แม้จะเขียนทีหลัง
}
```

## การเลือกใช้
1. **Channel**: 
   - ใช้สำหรับการสื่อสารระหว่าง goroutines
   - เหมาะกับการส่งข้อมูล
   - ใช้ง่าย ปลอดภัย

2. **Mutex**:
   - ใช้เมื่อต้องการล็อคการเข้าถึงข้อมูลที่ซับซ้อน
   - เหมาะกับการป้องกัน shared resources
   - ต้องระวังเรื่อง deadlock

3. **Atomic**:
   - ใช้กับ operations ง่ายๆ (เช่น +1, -1)
   - เร็วกว่า mutex
   - ใช้งานได้จำกัด

## ตัวอย่างใน goobserv
```go
// ใช้ channel สำหรับ logging
type Observer struct {
    buffer chan *Entry  // buffered channel
}

// Producer: ส่ง entry เข้า channel
obs.buffer <- entry

// Consumer: รับ entry จาก channel
for entry := range obs.buffer {
    output.Write(entry)
}
```

buffer chan *Entry เป็นการประกาศ channel ที่ใช้ส่ง pointer ของ Entry ครับ

แยกดูแต่ละส่วน:

chan คือ channel ใน Go - เป็นท่อส่งข้อมูลระหว่าง goroutines
*Entry คือ pointer ไปยัง Entry struct
buffer คือชื่อตัวแปร


// 1. สร้าง buffered channel ที่เก็บได้ 1000 entries
obs := NewObserver(&Config{
    BufferSize: 1000,
})
// ข้างในจะทำ: buffer: make(chan *Entry, 1000)

// 2. ส่ง entry เข้า channel
entry := NewEntry()
obs.buffer <- entry  // ส่ง entry เข้า channel

// 3. รับ entry จาก channel (ทำใน output)
for entry := range obs.buffer {
    // process entry
}


race condition คืออะไร
Race Condition คือสถานการณ์ที่โค้ดทำงานพร้อมกันหลาย goroutine และเข้าถึงข้อมูลเดียวกัน ทำให้เกิดผลลัพธ์ที่ไม่คาดคิด
ตัวอย่าง Race Condition:
var counter int = 0

```go
// Goroutine 1
go func() {
    counter++ // อ่าน counter = 0, บวก 1
}()

// Goroutine 2 (ทำพร้อมกัน)
go func() {
    counter++ // อ่าน counter = 0, บวก 1
}()

// ผลลัพธ์: counter = 1
// ที่ถูกควรเป็น: counter = 2
```

วิธีแก้ Race Condition:
1.ใช้ Mutex:
var counter int = 0
var mutex sync.Mutex

```go
go func() {
    mutex.Lock()
    counter++
    mutex.Unlock()
}()
```

2.ใช้ Channel:
```go
counter := make(chan int, 1)
counter <- 0 // ค่าเริ่มต้น

go func() {
    n := <-counter  // รับค่า
    counter <- n+1  // ส่งค่าใหม่
}()
```

03.ใช้ Atomic Operations:
```go
var counter atomic.Int64

go func() {
    counter.Add(1)
}()
```
ใน goobserv เราใช้ channel เพื่อป้องกัน race condition เพราะ:

ง่ายต่อการใช้งาน
ปลอดภัยโดยธรรมชาติ (thread-safe)
เหมาะกับการส่งข้อมูลระหว่าง goroutines

# Buffer ใน Go

## 1. โครงสร้างของ Buffer

Buffer ใน Go เป็น circular buffer (วงกลม) ที่มีโครงสร้างดังนี้:

```go
// สร้าง buffer ขนาด 1000
buffer := make(chan *Entry, 1000)

/*
โครงสร้างข้างใน:
[slot0][slot1][slot2]...[slot999]
 ^head            ^tail

- head: ตำแหน่งที่จะอ่านข้อมูลออก
- tail: ตำแหน่งที่จะเขียนข้อมูลเข้า
- เมื่อ head หรือ tail ถึงท้าย จะวนกลับมาที่ต้น
*/
```

## 2. การใช้งาน Buffer

### 2.1 Create (สร้าง)
```go
// แบบปกติ
buffer := make(chan *Entry, 1000)

// แบบมี config
type Observer struct {
    buffer chan *Entry
}

obs := NewObserver(&Config{
    BufferSize: 1000,
})
```

### 2.2 Write (เขียน)
```go
// เขียนข้อมูลเข้า buffer
buffer <- entry

// ถ้า buffer เต็ม จะ block จนกว่าจะมีที่ว่าง
// ควรใช้ select เพื่อป้องกัน block
select {
case buffer <- entry:
    // เขียนสำเร็จ
default:
    // buffer เต็ม, จัดการ error
}
```

### 2.3 Read (อ่าน)
```go
// อ่านทีละตัว
entry := <-buffer

// อ่านทั้งหมด
for entry := range buffer {
    // process entry
}
```

### 2.4 Clear (เคลียร์)
```go
// วิธีที่ 1: สร้างใหม่
buffer = make(chan *Entry, 1000)

// วิธีที่ 2: เคลียร์แบบ graceful
close(buffer)         // ปิด channel
for entry := range buffer {
    // process remaining entries
}
buffer = make(chan *Entry, 1000)  // สร้างใหม่
```

## 3. ทำไมใช้ Buffer Size 1000?

เป็นการป้องกันกรณีที่มี requests เข้ามาเยอะพร้อมกัน:

```
ตัวอย่างการคำนวณ:
- สมมติมี requests 2000 req/วินาที
- ระบบ process ได้ 500 entries/วินาที

ถ้า buffer น้อยเกินไป (10):
- buffer เต็มเร็ว
- requests ใหม่ต้องรอ (block)
- ระบบช้า

ถ้า buffer = 1000:
- รองรับ spike ได้
- requests ไม่ต้อง block
- output writer ทำงานทัน
```

## 4. การอัพเดท Span ใน Entry

สิ่งสำคัญ: ห้ามอัพเดท entry ที่อยู่ใน buffer!

วิธีที่ถูกต้อง:
```go
func Middleware(obs *Observer) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. สร้าง entry เก็บใน context
        entry := NewEntry()
        ctx := WithEntry(c.Request.Context(), entry)
        
        // 2. อัพเดท spans ระหว่าง request
        obs.Info(ctx, "Start")
        
        // ทำงานอื่นๆ...
        
        // 3. เมื่อ request เสร็จ
        entry.End()
        
        // 4. ส่ง entry ที่สมบูรณ์แล้วเข้า buffer
        obs.buffer <- entry
    }
}
```

สิ่งที่ไม่ควรทำ:
```go
// BAD: อัพเดท entry ที่อยู่ใน buffer
entry := NewEntry()
obs.buffer <- entry  // ส่งเข้า buffer ก่อน

entry.AddSpan(span)  // อันตราย! อาจเกิด race condition
```

เหตุผล:
1. ป้องกัน race condition
2. มั่นใจว่า entry มีข้อมูลครบก่อนส่ง
3. buffer ใช้ส่งต่อข้อมูลเท่านั้น ไม่ใช่ที่เก็บข้อมูล