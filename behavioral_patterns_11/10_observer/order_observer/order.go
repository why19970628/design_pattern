package main

import (
	"fmt"
	"sync"
)

// Subject 是主题接口，定义了添加、删除和通知观察者的方法。
type Subject interface {
	Attach(observer Observer, eventType string)
	Detach(observer Observer, eventType string)
	Notify(orderID string, status string)
}

// ConcreteSubject 是主题的具体实现。
type ConcreteSubject struct {
	observers map[string][]Observer
	mu        sync.Mutex
}

// NewConcreteSubject 创建一个新的 ConcreteSubject 实例。
func NewConcreteSubject() *ConcreteSubject {
	return &ConcreteSubject{
		observers: make(map[string][]Observer),
	}
}

// Attach 实现 Subject 接口中的 Attach 方法。
func (s *ConcreteSubject) Attach(observer Observer, eventType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers[eventType] = append(s.observers[eventType], observer)
}

// Detach 实现 Subject 接口中的 Detach 方法。
func (s *ConcreteSubject) Detach(observer Observer, eventType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if observers, ok := s.observers[eventType]; ok {
		for i, obs := range observers {
			if obs == observer {
				s.observers[eventType] = append(observers[:i], observers[i+1:]...)
				break
			}
		}
	}
}

// Notify 实现 Subject 接口中的 Notify 方法，通知对应事件类型的观察者。
func (s *ConcreteSubject) Notify(orderID string, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if observers, ok := s.observers[status]; ok {
		for _, observer := range observers {
			observer.Update(orderID)
		}
	}
}

// Observer 是观察者接口，定义了更新方法。
type Observer interface {
	Update(orderID string)
}

// ConcreteObserver 是观察者的具体实现。
type ConcreteObserver struct {
	name string
}

// Update 实现 Observer 接口中的 Update 方法。
func (o *ConcreteObserver) Update(orderID string) {
	fmt.Printf("Observer %s received order notification: Order ID %s\n", o.name, orderID)
	// 在这里可以添加具体的处理逻辑，比如调用仓库系统、物流系统、财务系统等进行相应的处理
}

func main() {
	// 创建主题对象
	subject := NewConcreteSubject()

	// 创建并注册观察者到主题中
	warehouseObserver := &ConcreteObserver{name: "Warehouse System"}
	logisticsObserver := &ConcreteObserver{name: "Logistics System"}
	financeObserver := &ConcreteObserver{name: "Finance System"}

	subject.Attach(warehouseObserver, "prepare")
	subject.Attach(logisticsObserver, "arrange")
	subject.Attach(financeObserver, "settle")

	// 模拟订单状态改变并通知对应的观察者
	orderID := "123456"
	subject.Notify(orderID, "prepare")
	subject.Notify(orderID, "arrange")
	subject.Notify(orderID, "settle")
}
