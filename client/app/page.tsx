"use client";

import { useState, useEffect } from "react";
import axios from "axios";
import styles from "../styles/Home.module.css";
import toast, { Toaster } from "react-hot-toast";

type Order = {
  id: string;
  item: string;
  amount: number;
};

export default function Home() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [newItem, setNewItem] = useState("");
  const [newAmount, setNewAmount] = useState(0);
  const [editingOrderId, setEditingOrderId] = useState<string | null>(null);
  const [editingItem, setEditingItem] = useState("");
  const [editingAmount, setEditingAmount] = useState(0);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    toast.loading("Fetching orders...");
    try {
      const res = await axios.get("http://localhost:8080/orders");
      setOrders(res.data);
      toast.dismiss();
      toast.success("Orders fetched successfully");
    } catch (error) {
      toast.dismiss();
      toast.error("Error fetching orders");
      console.error("Error fetching orders:", error);
    }
  };

  const createOrder = async () => {
    const newOrder = { item: newItem, amount: parseInt(newAmount.toString()) };
    toast.loading("Creating order...");
    try {
      await axios.post("http://localhost:8080/orders", newOrder, {
        headers: { "Content-Type": "application/json" },
      });
      fetchOrders();
      setNewItem("");
      setNewAmount(0);
      toast.dismiss();
      toast.success("Order created successfully");
    } catch (error) {
      toast.dismiss();
      toast.error("Error creating order");
      console.error("Error creating order:", error);
    }
  };

  const updateOrder = async (order: Order) => {
    toast.loading("Updating order...");
    try {
      await axios.put(
        `http://localhost:8080/orders/${order.id}/${order.item}`,
        order,
        {
          headers: { "Content-Type": "application/json" },
        }
      );
      fetchOrders(); // Refresh the order list
      setEditingOrderId(null);
      toast.dismiss();
      toast.success("Order updated successfully");
    } catch (error) {
      toast.dismiss();
      toast.error("Error updating order");
      console.error("Error updating order:", error);
    }
  };

  const deleteOrder = async (id: string, item: string) => {
    toast.loading("Deleting order...");
    try {
      await axios.delete(`http://localhost:8080/orders/${id}/${item}`);
      fetchOrders(); // Refresh the order list
      toast.dismiss();
      toast.success("Order deleted successfully");
    } catch (error) {
      toast.dismiss();
      toast.error("Error deleting order");
      console.error("Error deleting order:", error);
    }
  };

  const increaseAmount = async (id: string, item: string) => {
    const order = orders.find(
      (order) => order.id === id && order.item === item
    );
    if (order) {
      const updatedOrder = { ...order, amount: order.amount + 1 };
      toast.loading("Increasing amount...");
      try {
        await axios.put(
          `http://localhost:8080/orders/${id}/${item}`,
          updatedOrder,
          {
            headers: { "Content-Type": "application/json" },
          }
        );
        fetchOrders(); // Refresh the order list
        toast.dismiss();
        toast.success("Amount increased successfully");
      } catch (error) {
        toast.dismiss();
        toast.error("Error increasing amount");
        console.error("Error increasing amount:", error);
      }
    }
  };

  const decreaseAmount = async (id: string, item: string) => {
    const order = orders.find(
      (order) => order.id === id && order.item === item
    );
    if (order) {
      const updatedOrder = { ...order, amount: order.amount - 1 };
      toast.loading("Decreasing amount...");
      try {
        await axios.put(
          `http://localhost:8080/orders/${id}/${item}`,
          updatedOrder,
          {
            headers: { "Content-Type": "application/json" },
          }
        );
        fetchOrders(); // Refresh the order list
        toast.dismiss();
        toast.success("Amount decreased successfully");
      } catch (error) {
        toast.dismiss();
        toast.error("Error decreasing amount");
        console.error("Error decreasing amount:", error);
      }
    }
  };

  const startEditing = (order: Order) => {
    setEditingOrderId(order.id);
    setEditingItem(order.item);
    setEditingAmount(order.amount);
  };

  const cancelEditing = () => {
    setEditingOrderId(null);
  };

  const saveChanges = async () => {
    if (editingOrderId !== null) {
      const updatedOrder = {
        id: editingOrderId,
        item: editingItem,
        amount: editingAmount,
      };
      await updateOrder(updatedOrder);
    }
  };

  return (
    <div className={styles.container}>
      <Toaster position="top-center" reverseOrder={false} />
      <h1 className={styles.title}>Orders API</h1>
      <div className={styles.form}>
        <input
          className={styles.input}
          placeholder="Item"
          value={newItem}
          onChange={(e) => setNewItem(e.target.value)}
        />
        <input
          className={styles.input}
          placeholder="Amount"
          type="number"
          value={newAmount}
          onChange={(e) => setNewAmount(parseInt(e.target.value))}
        />
        <button className={styles.createButton} onClick={createOrder}>
          Create Order
        </button>
      </div>
      <h2 className={styles.subtitle}>Orders List</h2>
      <div className={styles.table}>
        {orders.map((order) => (
          <div key={`${order.id}-${order.item}`} className={styles.row}>
            <div className={styles.cell}>{order.id}</div>
            <div className={styles.cell}>
              {editingOrderId === order.id ? (
                <input
                  className={styles.input}
                  value={editingItem}
                  onChange={(e) => setEditingItem(e.target.value)}
                />
              ) : (
                order.item
              )}
            </div>
            <div className={styles.cell}>
              {editingOrderId === order.id ? (
                <input
                  className={styles.input}
                  type="number"
                  value={editingAmount}
                  onChange={(e) => setEditingAmount(parseInt(e.target.value))}
                />
              ) : (
                order.amount
              )}
            </div>

            {editingOrderId !== order.id && (
              <>
                <button
                  className={styles.increaseButton}
                  onClick={() => increaseAmount(order.id, order.item)}
                >
                  Increase Amount
                </button>
                <button
                  className={styles.decreaseButton}
                  onClick={() => decreaseAmount(order.id, order.item)}
                >
                  Decrease Amount
                </button>
              </>
            )}

            <button
              className={styles.deleteButton}
              onClick={() => deleteOrder(order.id, order.item)}
            >
              Delete
            </button>

            {editingOrderId === order.id ? (
              <>
                <button className={styles.updateButton} onClick={saveChanges}>
                  Save
                </button>
                <button className={styles.updateButton} onClick={cancelEditing}>
                  Cancel
                </button>
              </>
            ) : (
              <button
                className={styles.updateButton}
                onClick={() => startEditing(order)}
              >
                Edit
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
