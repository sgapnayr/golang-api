"use client";

import { useState, useEffect } from "react";
import axios from "axios";
import styles from "../styles/Home.module.css";

type Order = {
  id: string;
  item: string;
  amount: number;
};

export default function Home() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [item, setItem] = useState("");
  const [amount, setAmount] = useState(0);
  const [id, setId] = useState("");

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      const res = await axios.get("http://localhost:8080/orders");
      setOrders(res.data);
    } catch (error) {
      console.error("Error fetching orders:", error);
    }
  };

  const createOrder = async () => {
    const newOrder = { id, item, amount: parseInt(amount.toString()) };
    try {
      await axios.post("http://localhost:8080/orders", newOrder, {
        headers: { "Content-Type": "application/json" },
      });
      fetchOrders(); // Refresh the order list
      setItem("");
      setAmount(0);
      setId("");
    } catch (error) {
      console.error("Error creating order:", error);
    }
  };

  const updateOrder = async () => {
    const updatedOrder = { id, item, amount: parseInt(amount.toString()) };
    try {
      await axios.put(`http://localhost:8080/orders/${id}`, updatedOrder, {
        headers: { "Content-Type": "application/json" },
      });
      fetchOrders(); // Refresh the order list
    } catch (error) {
      console.error("Error updating order:", error);
    }
  };

  const deleteOrder = async (id: string) => {
    try {
      await axios.delete(`http://localhost:8080/orders/${id}`);
      fetchOrders(); // Refresh the order list
    } catch (error) {
      console.error("Error deleting order:", error);
    }
  };

  const increaseAmount = async (id: string, amount: number) => {
    try {
      await axios.put(
        `http://localhost:8080/orders/${id}/increase`,
        { amount },
        {
          headers: { "Content-Type": "application/json" },
        }
      );
      fetchOrders(); // Refresh the order list
    } catch (error) {
      console.error("Error increasing amount:", error);
    }
  };

  const decreaseAmount = async (id: string, amount: number) => {
    try {
      await axios.put(
        `http://localhost:8080/orders/${id}/decrease`,
        { amount },
        {
          headers: { "Content-Type": "application/json" },
        }
      );
      fetchOrders(); // Refresh the order list
    } catch (error) {
      console.error("Error decreasing amount:", error);
    }
  };

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Orders API</h1>
      <div className={styles.form}>
        <input
          className={styles.input}
          placeholder="ID"
          value={id}
          onChange={(e) => setId(e.target.value)}
        />
        <input
          className={styles.input}
          placeholder="Item"
          value={item}
          onChange={(e) => setItem(e.target.value)}
        />
        <input
          className={styles.input}
          placeholder="Amount"
          type="number"
          value={amount}
          onChange={(e) => setAmount(parseInt(e.target.value))}
        />
        <button className={styles.createButton} onClick={createOrder}>
          Create Order
        </button>
        <button className={styles.updateButton} onClick={updateOrder}>
          Update Order
        </button>
      </div>
      <h2 className={styles.subtitle}>Orders List</h2>
      <div className={styles.table}>
        {orders.map((order) => (
          <div key={order.id} className={styles.row}>
            <div className={styles.cell}>{order.id}</div>
            <div className={styles.cell}>{order.item}</div>
            <div className={styles.cell}>{order.amount}</div>
            <button
              className={styles.increaseButton}
              onClick={() => increaseAmount(order.id, 1)}
            >
              Increase Amount
            </button>
            <button
              className={styles.decreaseButton}
              onClick={() => decreaseAmount(order.id, 1)}
            >
              Decrease Amount
            </button>
            <button
              className={styles.deleteButton}
              onClick={() => deleteOrder(order.id)}
            >
              Delete
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}
