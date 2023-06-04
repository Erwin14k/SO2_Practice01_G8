// This library provides functions and structures necessary to create and manage Linux kernel modules.
#include <linux/module.h>
/* This library provides functions and macros for printing messages to the kernel console, including 
macros like printk() and KERN_INFO.*/
#include <linux/kernel.h>
/* This library defines the module_init() and module_exit() macros to register the module's 
initialization and exit function, respectively.*/
#include <linux/init.h>
/* This library provides functions and structures for creating and managing files in the procfs file 
system, which provides a user interface for accessing kernel information and statistics.*/
#include <linux/proc_fs.h>
// This library provides functions to copy data between user space and kernel space.
#include <asm/uaccess.h>
/* This library provides functions and structures for writing data to stream files. Stream files are an 
efficient way to write large amounts of data to a file without having to store everything in memory.*/	
#include <linux/seq_file.h>
// This library provides structures and functions for working with kernel processes and tasks.
#include <linux/sched.h>
/* This library provides functions and structures for working with the kernel's virtual memory space, 
such as the mm_struct structure, which describes the current state of a process's virtual memory space.*/
#include <linux/mm.h>
#include <linux/hugetlb.h>

// Module information
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("CPU Information Module");
MODULE_AUTHOR("Erwin Fernando Vasquez Peñate");

/* "cpu" is a pointer to a task_struct that will be used to loop through all running processes on the 
system.*/
struct task_struct* cpu;
// "child" is a pointer to a task_struct that will be used to traverse the child processes of a process.
struct task_struct* child;
/* lstProcess is a pointer to a data structure that is used to traverse the list of child processes
of a process.*/
struct list_head* lstProcess;
static int cpu_percentage(void);

//Function that will be executed every time the file is read with the CAT command
static int write_file(struct seq_file *the_file, void *v){   
    // Ram information
    struct sysinfo info;
    si_meminfo(&info);
    // Capture data in mb.
    // Total ram
    seq_printf(the_file, "{\"total_ram\":");
    seq_printf(the_file, "%ld", info.totalram* info.mem_unit / 1024/ 1024);
    // Free ram
    seq_printf(the_file, ",\"free_ram\":");
    seq_printf(the_file, "%ld", info.freeram* info.mem_unit / 1024 / 1024);
    // Occupied ram
    seq_printf(the_file, ",\"ram_occupied\":");
    seq_printf(the_file, "%ld", (info.totalram-info.freeram) * info.mem_unit / 1024 / 1024);
    seq_printf(the_file, "},\n");
    // CPU Information
    int ram, split, child_split;
    split = 0;
    child_split = 0;
    seq_printf(the_file, "[");
    for_each_process(cpu){
        if(split){
            seq_printf(the_file, ",");
        }
        seq_printf(the_file, "{\"pid\":");
        seq_printf(the_file, "%d", cpu->pid);
        seq_printf(the_file, ",\"name\":");
        seq_printf(the_file, "\"%s\"", cpu->comm);
        seq_printf(the_file, ",\"user\":");
        seq_printf(the_file, "%d", cpu->real_cred->uid);
        seq_printf(the_file, ",\"status\":");
        seq_printf(the_file, "%d", cpu->__state);
        if (cpu->mm) {
            ram = (get_mm_rss(cpu->mm)<<PAGE_SHIFT)/(1024*1024);
            seq_printf(the_file, ",\"ram\":");
            seq_printf(the_file, "%d", ram);
        }
        seq_printf(the_file, ",\"children\":[");
        child_split = 0;
        list_for_each(lstProcess, &(cpu->children)){
            child = list_entry(lstProcess, struct task_struct, sibling);
            if(child_split){
                seq_printf(the_file, ",");
            }
            seq_printf(the_file, "\n{ \"pid\" : %d, \"name\" : \"%s\"}", child->pid, child->comm);
            child_split = 1;
        }
        seq_printf(the_file, "]}\n");
        split = 1;
    }
    seq_printf(the_file, "]\n");
    return 0;
}

//Function that will be executed every time the file is read with the CAT command
static int when_open(struct inode *inode, struct file *file){
    return single_open(file, write_file, NULL);
}
//If the kernel is 5.6 or higher, use the proc_ops structure
static struct proc_ops operations ={
    .proc_open = when_open,
    .proc_read = seq_read
};

//Function to execute when inserting the module in the kernel with insmod
static int _insert(void){
    proc_create("mem_grupo8", 0, NULL, &operations);
    printk(KERN_INFO "Hola mundo, somos el grupo 8 y este es el monitor de memoria\n");
    return 0;
}
//Function to execute when removing the kernel module with rmmod
static void _remove(void){
    remove_proc_entry("mem_grupo8", NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 8 y este fue el monitor de memoria\n");
}
module_init(_insert);
module_exit(_remove);