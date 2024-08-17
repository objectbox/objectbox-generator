
#define OBX_CPP_FILE


#include "objectbox.hpp"

#include "objectbox-model.h"
#include "task.obx.hpp"
#include "monster.obx.hpp"
#include "another_monster.obx.hpp"
#include "person.obx.hpp"

#include <cinttypes>

int main(int argc, char* args[])
{
    // create_obx_model() provided by objectbox-model.h
    // obx interface contents provided by objectbox.hpp
    obx::Store store(create_obx_model());
    obx::Box<Task> box(store);

    Task my_task{};
    my_task.text = "Buy milk";
    obx_id id = box.put(my_task);  // Create

    std::unique_ptr<Task> task = box.get(id);   // Read
    if (task) {
        task->text += " & some bread";
        box.put(*task);                         // Update
    }

    printf("Your task has ID=%" PRIu64 ", text=%s\n",
        id,
        box.get(id)->text.c_str());

    return 0;
}
