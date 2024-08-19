
#define OBX_CPP_FILE


#include "objectbox.hpp"

#include "objectbox-model.h"
#include "task.obx.hpp"
#include "monster.obx.hpp"
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

    printf("Your task has ID=%" PRIu64 ", text=%s\n, text2=%s\n",
        id,
        box.get(id)->text.c_str(),
        box.get(id)->text2.c_str());
    
    obx::Box<Person> person_box(store);
    Person my_person{};
    obx_id person_id = person_box.put(my_person);  // Create
    
    printf("Your person has ID=%" PRIu64, person_id);

    return 0;
}
