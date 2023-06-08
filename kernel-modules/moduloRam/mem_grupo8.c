//Header to use the memory info
#include <linux/mm.h> 
//Allows access to functions and structures necessary for the creation and administration of kernel modules. (Mandatory)
#include <linux/module.h>
// To use kern info
#include <linux/kernel.h>
//Header to use the init and exit macros
#include <linux/init.h>
//Header necesary to use proc_ops/file_operations
#include <linux/proc_fs.h>
/* for copy_from_user */
#include <asm/uaccess.h>	
/* Header to use the seq_file library and manage the file on  /proc*/
#include <linux/seq_file.h>


// Module info
MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo RAM, Laboratorio Sistemas Operativos 2");
MODULE_AUTHOR("GRUPO8");

// Function to write the file
static int escribir_archivo(struct seq_file *archivo, void *v)
{
    // struct sysinfo
    struct sysinfo info;
    si_meminfo(&info);

    //Init the file with a {
    seq_printf(archivo, "{\n");
    //The total ram on the system
    seq_printf(archivo, "\"totalram\":%lu,\n", info.totalram * info.mem_unit / 1024 / 1024);
    //The free ram on the system
    seq_printf(archivo, "\"ramlibre\":%lu,\n", info.freeram * info.mem_unit / 1024 / 1024);
    //The occupied ram on the system
    seq_printf(archivo, "\"ramocupada\":%lu\n", (info.totalram -  info.freeram )* info.mem_unit / 1024 / 1024);
    //Closing the file with a }
    seq_printf(archivo, "}\n");

    return 0;
}

//Funcion executed every time with a cat command 
static int al_abrir(struct inode *inode, struct file *file)
{
    return single_open(file, escribir_archivo, NULL);
}

// If the kernel is 5.6 or higher
static struct proc_ops operaciones =
{
    .proc_open = al_abrir,
    .proc_read = seq_read
};


// Function to display the message when the module is created with insmod
static int _insert(void)
{
    proc_create("mem_grupo8", 0, NULL, &operaciones);
    printk(KERN_INFO "Hola mundo, somos el grupo 8 y este es el monitor de memoria\n");
    return 0;
}
// Function to display the message when the module is deleted with rmmod
static void _remove(void)
{
    remove_proc_entry("mem_grupo8", NULL);
    printk(KERN_INFO "Sayonara mundo, somos el grupo 8 y este fue el monitor de memoria\n");
}

module_init(_insert);
module_exit(_remove);