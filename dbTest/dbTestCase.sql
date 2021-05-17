-- a. Display Customer List including calculating the total order.
SELECT c.customer_id, c.customer_name, count(o2.customer_id) as total_order from customers c, orders o
where o.customer_id = c.customer_id group by c.customer_id, o.customer_id

-- b. Show Product List including calculating the number of orders sorted by the most in the order.
SELECT p.product_name, p.basic_price, p.product_name, count(o.order_detail_id) 
as numberOfOrder from products p, orderdetails o where p.product_id = o.product_id 
group by p.product_id, o.product_id order by count(o.order_detail_id) desc

-- c. Display the sort payment method data most frequently used by customers.
select p.method_name, p.payment_method_id, count(o.payment_method_id) 
as TotalUsed from paymentmethods p, orders o 
where p.payment_method_id = o.payment_method_id
group by p.method_name, p.payment_method_id 
order by count(o.payment_method_id) desc