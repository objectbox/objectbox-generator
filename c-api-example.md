ObjectBox Generator: C example
==============================

Running `objectbox-generator -c tasklist.fbs` will generate C binding code for 
`tasklist.fbs` - we get the following files:

* objectbox-model.h
* objectbox-model.json
* tasklist.obx.h

> Note: you should add all these files to your source control (e.g. git), 
> most importantly the objectbox-model.json which ensures compatibility 
> with previous versions of your database after you make changes to the schema.

Now in your application, you can include the headers and start to work with your database. 
Have a look at the following `main.c` showing one of the many ways you can work with 
objectbox-c and the generated code:

```c
#include "objectbox.h"
#include "objectbox-model.h"
#include "tasklist.obx.h"

obx_err print_last_error() {
    printf("Unexpected error: %d %s\n", 
        obx_last_error_code(), obx_last_error_message());
    return obx_last_error_code();
}

obx_id task_put(OBX_box* box, Task* task) {
    flatcc_builder_t builder;
    flatcc_builder_init(&builder);

    size_t size = 0;
    void* buffer = NULL;

    // Note: Task_to_flatbuffer() is provided by the generated code
    obx_id id = 0;
    if (Task_to_flatbuffer(&builder, task, &buffer, &size)) {
        id = obx_box_put_object(box, buffer, size, OBXPutMode_PUT);  // 0 on error
    }

    flatcc_builder_clear(&builder);
    if (buffer) flatcc_builder_aligned_free(buffer);

    if (id == 0) {
        // TODO: restructure; won't print the right error if it occurred 
        //  in Task_to_flatbuffer(), i.e. outside objectbox
        print_last_error();
    } else {
        task->id = id;  // update the ID property on new objects for convenience
    }

    return id;
}

Task* task_read(OBX_store* store, OBX_box* box, obx_id id) {
    OBX_txn* txn = NULL;

    // We need an explicit TX - read data lifecycle is bound to the open TX.
    // The transaction can be closed after reading the object from flatbuffers.
    txn = obx_txn_read(store);
    if (!txn) {
        print_last_error();
        return NULL;
    }

    void* data;
    size_t size;
    int rc = obx_box_get(box, id, &data, &size);
    if (rc != OBX_SUCCESS) {
        // if (rc == OBX_NOT_FOUND); // No special treatment at the moment
        obx_txn_close(txn);
        return NULL;
    }

    Task* result = Task_new_from_flatbuffer(data, size);
    obx_txn_close(txn);
    return result;
}

int main(int argc, char* args[]) {
    int rc = 0;
    OBX_store* store = NULL;
    OBX_box* box = NULL;
    Task* task = NULL;

    // Firstly, we need to create a model for our data and the store
    {
        OBX_model* model = create_obx_model();  // defined in objectbox-model.h
        if (!model) goto handle_error;
        if (obx_model_error_code(model)) {
            printf("Model definition error: %d %s\n", 
                obx_model_error_code(model), obx_model_error_message(model));
            obx_model_free(model);
            goto handle_error;
        }

        OBX_store_options* opt = obx_opt();
        obx_opt_model(opt, model);
        store = obx_store_open(opt);
        if (!store) goto handle_error;

        // obx_store_open() takes ownership of model and opt and frees them.
    }

    box = obx_box(store, Task_ENTITY_ID);  // Note the generated "Task_ENTITY_ID"

    obx_id id = 0;

    {  // Create
        Task task = {.text = "Buy milk"};
        id = task_put(box, &task);
        if (!id) goto handle_error;
        printf("New task inserted with ID %d\n", id);
    }

    {  // Read
        task = task_read(store, box, id);
        if (!task) goto handle_error;
        printf("Task %d read with text: %s\n", id, task->text);
    }

    {  // Update
        const char* appendix = " & some bread";

        // updating a string property is a little more involved 
        size_t old_text_len = task->text ? strlen(task->text) : 0;
        char* new_text = 
            (char*) malloc((old_text_len + strlen(appendix) + 1) * sizeof(char));

        if (task->text) {
            memcpy(new_text, task->text, old_text_len);

            // free the memory allocated previously before overwritting below
            free(task->text);
        }
        memcpy(new_text + old_text_len, appendix, strlen(appendix) + 1);
        task->text = new_text;
        printf("Updated task %d with a new text: %s\n", id, task->text);
    }

    // Delete
    if (obx_box_remove(box, id) != OBX_SUCCESS) goto handle_error;

free_resources:  // free any remaining allocated resources
    if (task) Task_free(&task); // free allocs by Task_new_from_flatbuffer()
    if (store) obx_store_close(store); // and close the store
    return rc;

handle_error:  // print error and clean up
    rc = print_last_error();
    if (rc <= 0) rc = 1;
    goto free_resources;
}
```

To compile, link to the objectbox-c library and flatcc-runtime library, 
e.g. something like this should work: `gcc main.c -I. -lobjectbox -lflatccrt`. 
Note: the command snippet assumes you have objectbox-c and flatccrt libraries installed in a path 
recognized by your OS (e.g. /usr/local/lib/) and all the referenced headers are in the same folder as `main.c`.
