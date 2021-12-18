### Properties of B Tree
- k >= 1
- each leaf is at a same distance from the root
- each node holds upto 2*k keys and at least upto k keys except root
- root can a minimum of 1 key and max of 2*k keys
### type BTree
- Contains a field called as k
- Less Function for comparing the keys
- Method Insert(k Key)
    - child <- root
    - current <- child
    - stack <- []Page{current}
    - u <- Entry{key: k}
    - loop isLeaf(child)
    - child <- scan(current, u); stack.push(child)
    - child = stack.pop()
    - label:
        - insert(current, u)
        - if !IsSafe(current)
            - middleEntry, pageRight := SplitMiddle(current)
            - pageRight.head.entry.pagePtr = middleEntry.pagePtr
            - middleEntry.pagePtr = pageRight
            - u <- middleEntry
            - current <- stack.pop
            - if current is nil
                - page := Newpage()
                - insert(page, u)
                - root <- page
            - else goto label
    - else
        - stack.push(pos)
### type Page
- contains a doubly linked list so that it can be useful to find both right and left siblings
- method to add a new Entry returns the same Entry
- contains a field called as length 
- As the adding always happens on the leaf node
- Scan
    - returns the first Entry's page pointer which is just greater or than or equal
    - If matched Entry's pagePtr is nil then return current
- Insert
    - inserts an Entry at pos which is just greater than or equal 
- Split Middle
    - creates a new Page with mid+1..n entries in it
    - removes mid+1..n entries in it
    - return mid Entry and new Page ptr
### type Entry
- contains a key and the page ptr key
### type interface Key
- Implements Less(than Key) bool
### type Stack arr of Page
- pop
- push
