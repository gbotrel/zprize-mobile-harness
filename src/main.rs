use celo_zprize::{
    benchmark_msm, deserialize_input, gen_random_vectors, gen_zero_vectors, serialize_input,
    FileInputIterator,
};
use rand::thread_rng;

const INSTANCE_SIZE: usize = 16;
const NUM_INSTANCES: usize = 10;

fn main() {
    let dir = format!("./vectors/{}x{}", NUM_INSTANCES, INSTANCE_SIZE);
    //let dir = "android/ZPrize/app/src/main/assets/".to_string();
    //generate_vectors(&dir);
    run_benchmark(&dir);
}

fn generate_vectors(dir: &str) {
    let mut rng = thread_rng();
    println!("Generating elements");
    let n_elems = 1 << INSTANCE_SIZE;
    for i in 0..NUM_INSTANCES {
        let (points, scalars) = gen_random_vectors(n_elems, &mut rng);
        serialize_input(dir, &points, &scalars, i != 0).unwrap();
    }
    println!("Generated elements");
}

fn run_benchmark(dir: &str) {
    println!("Running benchmark for baseline result");
    let input_iter = FileInputIterator::open(dir).unwrap();
    let result = benchmark_msm(dir, input_iter, 1);
    println!("Done running benchmark: {:?}", result);
}
